### Software utilized
- Golang
- Docker
- RabbitMQ
- Kubernetes
- PostgreSQL

### Core Features
- [x] Task Scheduling and Execution:
    - Basic task queue where tasks can be enqueued and dequeued.
    - Asynchronous task execution with workers.

- [ ] Distributed Architecture:
    - Distributed system design with multiple worker nodes.
    - Tasks are distributed across nodes.

- [ ] Fault Tolerance:
    - Retries for failed tasks.
    - Heartbeats to detect worker failures and reassign tasks.

- [ ] Scalability:
    - Can scale horizontally by adding more worker nodes.
    - Load balancer / consistent hashing to distribute tasks evenly.

### Nice Features
- [ ] Priority Queues:
    - Support for task prioritization (high, medium, low priority).

- [ ] Task Dependencies:
    - Tasks can have dependencies (e.g., Task B runs only after Task A completes).

- [ ] Rate Limiting:
    - Rate limiting for tasks to prevent overloading the system.

- [ ] Task Timeouts:
    - Timeouts for tasks to ensure they don’t run indefinitely.

- [ ] Persistence:
    - PostgreSQL to persist tasks and their states.
    - Tasks can be recovered after a system crash.

- [ ] Distributed Consensus:
    - Consensus algorithm like Raft or leverage etcd/Zookeeper for coordination between nodes.

- [ ] Monitoring and Metrics:
    - Prometheus / OpenTelemetry to collect metrics (e.g., task completion rate, worker load).

- [ ] Dashboard to visualize system performance.
    - Grafana

- [ ] Dynamic Worker Scaling:
    - Implement auto-scaling for workers based on task load by using Kubernetes

- [x] Task Batching:
    - Batching of small tasks to improve efficiency.

- [ ] Dead Letter Queue:
    - Dead letter queue for tasks that fail repeatedly.

- [ ] API and CLI:
    - [x] A REST/gRPC API for enqueuing tasks and checking their status.
    - Build a CLI tool for interacting with the task queue.

- [x] Event-Driven Architecture:
    - Uses RabbitMQ message broker for task distribution.

- [ ] Distributed Tracing:
    - Integrate distributed tracing  by Jaeger to track tasks across workers.

- [ ] Task Result Storage:
    - Store task results in a distributed storage system in S3.

- [ ] Graceful Shutdown:
    - Graceful shutdown for workers to complete ongoing tasks before exiting.

- [ ] Custom Task Routing:
    - Tasks to be routed to specific workers based on metadata or tags.

- [ ] Idempotency:
    - Ensure tasks are idempotent to handle duplicate executions gracefully.

- [ ] Security:
    - Authentication and authorization for API access.
    - Encrypt sensitive task data.

- [ ] Testing and CI/CD:
    - Unit tests.
    - Integration tests.
    - Load tests.
    - CI/CD pipeline for automated testing and deployment.

### How to run
```docker compose up --build -d``` and see it in action
It has live build, docker no need to restart Docker every time new codes are added or updated.
