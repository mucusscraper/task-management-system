# Task Management API

A robust, enterprise-grade Task Management System built in Go (Golang) using the **Gin Gonic** web framework, **PostgreSQL** as the relational database, and **Goose** for version-controlled database migrations. 

This API models a secure, forward-only task lifecycle managed by two user roles: **Supervisors** and **Workers**.

---

## 🏗️ Layered Architecture & Design

The project is structured following clean, layered architecture principles to enforce separation of concerns, high testability, and decoupling of business logic from HTTP transport layers.

```text
.
├── cmd/
│   └── server/
│       └── main.go                 # App entrypoint (initializes DB, router, and starts server)
├── internal/
│   ├── database/
│   │   └── postgres.go             # SQL driver connection & embedded Goose migration engine
│   ├── dto/
│   │   └── task.go                 # Strongly-typed Request/Response Data Transfer Objects
│   ├── handlers/
│   │   ├── methods.go              # HTTP controller handlers (reads context, returns JSON/Status codes)
│   │   ├── middleware.go           # Authentication & Context-setting Middleware
│   │   ├── middleware_test.go      # Middleware unit tests
│   │   └── task_handler.go         # Handler struct mapping
│   ├── models/
│   │   ├── notification.go         # Domain model for notification log schema
│   │   ├── task.go                 # Domain model for tasks & status types
│   │   └── user.go                 # Domain model for users & role types
│   └── service/
│       ├── errors.go               # Centralized, domain-specific sentinel errors
│       ├── methods.go              # Core Business Logic Layer (enforces state machine & rules)
│       ├── methods_test.go         # Complete service & DB integration tests
│       └── task_service.go         # Service struct mapping
└── migrations/                     # Raw SQL database schema migrations managed by Goose
```

## Flow of a Request
```
[HTTP Client] ──► [Auth Middleware] ──► [Task Handlers] ──► [Task Service] ──► [PostgreSQL]
 (Headers)      (Injects User Context)    (Payload DTO)      (Business Rules)   (Data Persistence)
```

 ## Task Status Lifecycle & Business Rules

 The application implements a strict, forward-only finite state machine to manage task progression safely:
 ```
 [CREATED] ──(AssignTask)──► [ASSIGNED] ──(UpdateStatus)──► [IN_PROGRESS] ──(UpdateStatus)──► [COMPLETED]
```

## Security
- Task Creation: Only Supervisors can create a task. Newly created tasks always start with the ```CREATED``` status.

- Task Assignment: Only Supervisors can assign a task to a Worker. This transitions the task to ```ASSIGNED```. When assigned, the system logs a terminal message and persists a physical notification record in the database.

- Status Transitions: Transitions are forward-only. Backward transitions are blocked.

- Ownership Enforcement: A Worker can only view tasks assigned to them, and only the assigned Worker can transition their tasks to ```IN_PROGRESS``` or ```COMPLETED```. Any attempt by an unassigned worker to modify or fetch the task returns a ```403 Forbidden error```.

## Prerequisites & Installation

1. Database configuration:
The application is configured to connect to PostgreSQL using the following credentials by default: 

- Host: ```localhost```
- Port: ```5432```
- User: ```postgres```
- Password: ```postgres```
- Database Name: ```task_management```

Ensure you have created the database before starting the application:
```
CREATE DATABASE task_management;
```

2. Auto-Migrations & Seeds:
The application utilizes an embedded Goose migration engine. On startup, it automatically runs all SQL files in the ```migrations/``` directory, creating the schema and seeding the database with the required mock users.

| ID | Name       | Role       | HTTP Header for Authentication |
| -- | ---------- | ---------- | ------------------------------ |
| 1  | Supervisor | SUPERVISOR | X-User-Id: 1                   |  
| 2  | Worker 1   | WORKER     | X-User-Id: 2                   | 
| 3  | Worker 2   | WORKER     | X-User-Id: 3                   |

## Running the API and Tests

To start the HTTP server on port ```:8080```, run:
```
go run cmd/server/main.go
```

To run the automated test suite (includes both middleware unit tests and complete database-level integration tests), run:
```
go test ./... -v
```

## API Endpoints & cURL Testing Guide

Open a secondary terminal to test the complete, secure task lifecycle step-by-step:

1. System Health Check (Stretch Goal)
Validate that the API is active and successfully connected to the PostgreSQL instance:
```
curl -i -X GET http://localhost:8080/health
```
Expected Response: ```200 OK``` with ```{"database":"connected","status":"healthy"}```

2. Create Task (Supervisor Only)
Create a task as a Supervisor (User 1):
```
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{"title": "Deploy System", "description": "Deploying application to production server"}'
```
Expected Response: ```201 Created``` with task details ```(Status: CREATED, ID: 1)```.

Try creating a task as a Worker (User 2) to test unauthorized rejection:
```
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 2" \
  -d '{"title": "Unauthorized Task", "description": "Should fail"}'
```
Expected Response: ```403 Forbidden``` with ```{"error": "only supervisors can perform this action"}```.

3. Assign Task (Supervisor Only)
Assign the task (ID 1) to Worker 1 (User 2):
```
curl -i -X POST http://localhost:8080/tasks/1/assign \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{"worker_id": 2}'
```
Expected Response: ```200 OK``` with status changed to ```ASSIGNED```. The console running the server will print: ```Task 1 assigned to worker 2```.

4. Update Task Status (Assigned Worker Only)
The assigned worker Worker 1 (User 2) sets the task to ```IN_PROGRESS```:
```
curl -i -X PATCH http://localhost:8080/tasks/1/status \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 2" \
  -d '{"status": "IN_PROGRESS"}'
```
Expected Response: ```200 OK``` with status ```IN_PROGRESS```.

Try updating this task as an unassigned worker Worker 2 (User 3) to check security rules:
```
curl -i -X PATCH http://localhost:8080/tasks/1/status \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 3" \
  -d '{"status": "COMPLETED"}'
```
Expected Response: ```403 Forbidden``` with ```{"error": "only the worker of this task can perform this action"}```

5. Fetch Role-Based Task Lists
- As Supervisor (Retrieves all tasks in the system):
```
curl -i -X GET http://localhost:8080/tasks -H "X-User-Id: 1"
```
- As Worker 1 (Filters and retrieves only tasks assigned to them):
```
curl -i -X GET http://localhost:8080/tasks -H "X-User-Id: 2"
```

6. Get Worker Notifications (Stretch Goal)
Retrieve a historical record of all notifications associated with Worker 1 (User 2):
```
curl -i -X GET http://localhost:8080/tasks/notifications -H "X-User-Id: 2"
```
Expected Response: ```200 OK``` with a JSON array containing stored notification records.

## Reflection & Architecture Review
### How did you model tasks and status?
- Models: Relational tables were designed to enforce data integrity. The ```tasks``` table uses an explicit PostgreSQL ```CHECK``` constraint ensuring that the status column can only contain ```'CREATED'```, ```'ASSIGNED'```, ```'IN_PROGRESS'```, or ```'COMPLETED'```.

- State Safety: Transitions are validated programmatically inside the ```Service``` layer using domain-driven sentinel errors (e.g., ```ErrInvalidTransition```). This acts as an API safeguard before hitting the database.

### How do you enforce who can do what?
- Authentication: A lightweight ```AuthMiddleware``` parses the ```X-User-Id``` request header, queries the database to verify user existence, and injects the retrieved ```models.User``` object directly into the Gin Context (```ctx.Set("currentUser", user```)).

- **Authorization**: The Service layer extracts the injected user role. Supervisor-only or Worker-only capabilities are strictly checked through conditional blocks, returning custom error types that are transparently converted into ```HTTP 403 Forbidden``` or ```400 Bad Request``` codes by a centralized handler helper (```handleError```).

### What would you add with more time?
- **Database Transactions**: Grouping the task update/assignment operations and the notification insertions inside a safe ```db.BeginTx()``` block to roll back state changes if the database notification write fails.