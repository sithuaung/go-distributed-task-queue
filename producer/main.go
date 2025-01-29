package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	h "github.com/sithuaung/go-distributed-task-queue/helpers"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
)

// Task struct for tasks sent to RabbitMQ
type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority" validate:"oneof=high medium low"`
}

// RabbitMQ connection struct
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
	bq   amqp.Queue
}

// Global RabbitMQ instance
var rabbit *RabbitMQ

var Priorities = map[string]int{
	"high":   10,
	"medium": 5,
	"low":    1,
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection
func initRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{
			"x-max-priority": int32(10), // Maximum priority level
		},
	)
	failOnError(err, "Failed to declare a queue")

	bq, err := ch.QueueDeclare(
		"batch_task_queue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		amqp.Table{
			"x-max-priority": int32(10), // Maximum priority level
		},
	)
	failOnError(err, "Failed to declare a queue")

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
		q:    q,
		bq:   bq,
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task format", http.StatusBadRequest)
		return
	}

	err := validate.Struct(task)
	if err != nil {
		fmt.Printf("Task '%s' has invalid priority: %s\n", task.Title, err)
	}

	task.ID = uuid.New().String()

	// Serialize the task into JSON
	body, err := json.Marshal(&task)
	failOnError(err, "Failed to serialize Task")

	// Publish the task
	err = rabbit.ch.Publish(
		"",            // exchange
		rabbit.q.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "application/json",
			Body:         []byte(body),
			Priority:     uint8(Priorities[task.Priority]),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent task: %v\n", task)
}

func createBatchTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	tasks, err := validateTasks(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid payload: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println("Tasks: ", tasks)

	// Serialize the task into JSON
	batch, err := json.Marshal(tasks)
	failOnError(err, "Failed to serialize Tasks")

	// Publish the task
	err = rabbit.ch.Publish(
		"",             // exchange
		rabbit.bq.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "application/json",
			Body:         batch,
			Priority:     uint8(calculateBatchPriority(tasks)),
		})

	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent batch task: %v\n", tasks)
}

func main() {
	ctx := context.Background()

	// Initialize OpenTelemetry
	tp, mp, lp, err := h.InitOpenTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer h.ShutdownOpenTelemetry(ctx, tp, mp, lp)

	// Create a tracer
	tracer := otel.Tracer("rabbitmq-producer")

	// Start a span
	ctx, span := tracer.Start(ctx, "produce-tasks")
	defer span.End()

	rabbit = initRabbitMQ()
	defer rabbit.conn.Close()
	defer rabbit.ch.Close()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTaskHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/batch-tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createBatchTaskHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":8080"
	log.Printf("API server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func validateTasks(payload []byte) ([]Task, error) {
	// Step 1: Attempt to unmarshal the payload into a slice of Task
	var tasks []Task
	if err := json.Unmarshal(payload, &tasks); err != nil {
		return nil, errors.New("invalid JSON format or not an array of tasks")
	}

	// Step 2: Validate each task
	for i, task := range tasks {

		if task.Title == "" {
			return nil, errors.New("each task must have a title")
		}

		if task.Priority != "high" && task.Priority != "medium" && task.Priority != "low" {
			return nil, fmt.Errorf(
				"invalid priority for task '%s': must be 'high', 'medium', or 'low'",
				task.Title,
			)
		}

		tasks[i].ID = uuid.New().String()
	}

	return tasks, nil
}

func calculateBatchPriority(tasks []Task) int {
	priorityMap := map[string]int{
		"high":   10,
		"medium": 5,
		"low":    1,
	}
	highest := 0
	for _, task := range tasks {
		if p, ok := priorityMap[task.Priority]; ok && p > highest {
			highest = p
		}
	}
	return highest
}
