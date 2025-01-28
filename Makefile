send_task:
	curl -X POST \
		-H "Content-Type: application/json" \
		-d '{"id": "some-id", "title":"Example title", "priority": "low"}' \
		http://localhost:8080/tasks

send_batch_tasks:
	curl -X POST http://localhost:8080/batch-tasks \
		-H "Content-Type: application/json" \
		-d '[{"title":"Task 1","priority":"high"},{"title":"Task 2","priority":"medium"},{"title":"Task 3","priority":"low"}]'

up:
	docker compose up -d --build
