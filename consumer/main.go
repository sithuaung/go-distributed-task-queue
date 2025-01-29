package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/streadway/amqp"

	"github.com/sithuaung/go-distributed-task-queue/otel"
)

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority,omitempty"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func processTask(task Task, taskType string) {
	fmt.Printf(
		"Processing task: id:%v, title:%v, priority:%v\n",
		task.ID,
		task.Title,
		task.Priority,
	)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	fmt.Printf("Task [%v] processing complete!", task.ID)
}

func main() {
	otel.InitTracer()

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{
			"x-max-priority": int32(10), // Same max priority as the producer
		},
	)
	failOnError(err, "Failed to declare a queue")

	// Declare a queue
	bq, err := ch.QueueDeclare(
		"batch_task_queue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		amqp.Table{
			"x-max-priority": int32(10), // Same max priority as the producer
		},
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (set to false for manual acknowledgment)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register msgs consumer")

	// Consume messages
	batch_msgs, err := ch.Consume(
		bq.Name, // queue
		"",      // consumer
		false,   // auto-ack (set to false for manual acknowledgment)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register batch_msgs consumer")

	// Process messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Deserialize the batch
			var task Task
			err := json.Unmarshal(d.Body, &task)
			failOnError(err, "Failed to deserialize task")

			// Process the task
			processTask(task, "single")

			// Acknowledge the message
			d.Ack(false)
		}
	}()

	go func() {
		for d := range batch_msgs {
			// Deserialize the batch
			var tasks []Task
			err := json.Unmarshal(d.Body, &tasks)
			failOnError(err, "Failed to deserialize batch")

			// Sort tasks by priority
			sort.SliceStable(tasks, func(i, j int) bool {
				priorityMap := map[string]int{"high": 3, "medium": 2, "low": 1}
				return priorityMap[tasks[i].Priority] > priorityMap[tasks[j].Priority]
			})

			// Process the batch
			for _, task := range tasks {
				processTask(task, "batch")
			}

			// Acknowledge the message
			d.Ack(false)
		}
	}()

	fmt.Println(" [*] Waiting for tasks. To exit, press CTRL+C")
	<-forever
}
