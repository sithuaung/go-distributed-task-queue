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

3. Fault Tolerance:
    - Retries for failed tasks.
    - Heartbeats to detect worker failures and reassign tasks.

4. Scalability:
    - Can scale horizontally by adding more worker nodes.
    - Load balancer / consistent hashing to distribute tasks evenly.

### Nice Features
1. Priority Queues:
    - Support for task prioritization (high, medium, low priority).

2. Task Dependencies:
    - Tasks can have dependencies (e.g., Task B runs only after Task A completes).

3. Rate Limiting:
    - Rate limiting for tasks to prevent overloading the system.

4. Task Timeouts:
    - Timeouts for tasks to ensure they don’t run indefinitely.

5. Persistence:
    - PostgreSQL to persist tasks and their states.
    - Tasks can be recovered after a system crash.

6. Distributed Consensus:
    - Consensus algorithm like Raft or leverage etcd/Zookeeper for coordination between nodes.

7. Monitoring and Metrics:
    - Prometheus / OpenTelemetry to collect metrics (e.g., task completion rate, worker load).

8. Dashboard to visualize system performance.
    - Grafana

9. Dynamic Worker Scaling:
    - Implement auto-scaling for workers based on task load by using Kubernetes

10. Task Batching:
    - Batching of small tasks to improve efficiency.

11. Dead Letter Queue:
    - Dead letter queue for tasks that fail repeatedly.

12. API and CLI:
    - A REST/gRPC API for enqueuing tasks and checking their status.
    - Build a CLI tool for interacting with the task queue.

13. Event-Driven Architecture:
    - Uses RabbitMQ message broker for task distribution.

14. Distributed Tracing:
    - Integrate distributed tracing  by Jaeger to track tasks across workers.

15. Task Result Storage:
    - Store task results in a distributed storage system in S3.

16. Graceful Shutdown:
    - Graceful shutdown for workers to complete ongoing tasks before exiting.

15. Custom Task Routing:
    - Tasks to be routed to specific workers based on metadata or tags.

16. Idempotency:
    - Ensure tasks are idempotent to handle duplicate executions gracefully.

17. Security:
    - Authentication and authorization for API access.
    - Encrypt sensitive task data.

18. Testing and CI/CD:
    - Unit tests.
    - Integration tests.
    - Load tests.
    - CI/CD pipeline for automated testing and deployment.

### How to run
```docker compose up --build -d``` and see it in action
It has live build, docker no need to restart Docker every time new codes are added or updated.
