# SSUbench

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

A REST API platform where customers post tasks, executors respond to them, and completed work is paid with virtual points.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [Makefile Commands](#makefile-commands)
- [API Reference](#api-reference)
- [Roles](#roles)
- [Task Statuses](#task-statuses)
- [curl Examples](#curl-examples)
- [Tests](#tests)
- [License](#license)

---

## Tech Stack

- **Go 1.23** — primary language
- **PostgreSQL 16** — database
- **chi** — HTTP router
- **pgx/v5** — PostgreSQL driver
- **golang-migrate** — database migrations
- **JWT** — authentication
- **bcrypt** — password hashing
- **Docker / docker-compose** — infrastructure

---

## Architecture
```
cmd/api/          — entry point
internal/
  config/         — configuration
  domain/         — entities, constants, errors
  repo/           — database layer (interfaces + implementations)
  service/        — business logic
  handler/        — HTTP handlers and router
  middleware/     — request_id, logger, recover, auth
migrations/       — SQL migration files
docs/             — openapi.yaml
```

Layers communicate only through interfaces: `handler → service → repo`.

---

## Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/FelixWinchester/ssubench.git
cd ssubench
```

### 2. Create `.env`
```bash
cp .env.example .env
```

Edit `.env` if needed. Default values work out of the box for local development.

### 3. Start PostgreSQL
```bash
make docker-up
```

### 4. Run migrations
```bash
make migrate-up
```

### 5. Start the server
```bash
make run
```

Server will be available at `http://localhost:8080`.

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `SERVER_READ_TIMEOUT` | Read timeout | `10s` |
| `SERVER_WRITE_TIMEOUT` | Write timeout | `10s` |
| `SERVER_IDLE_TIMEOUT` | Idle timeout | `60s` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `ssubench` |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `JWT_SECRET` | JWT signing secret | — |
| `JWT_TTL` | Token lifetime | `24h` |
| `LOG_LEVEL` | Log level | `debug` |

---

## Makefile Commands
```bash
make run          # run the server
make build        # build binary to bin/api
make test         # run tests
make lint         # run linter
make docker-up    # start PostgreSQL
make docker-down  # stop PostgreSQL
make migrate-up   # apply migrations
make migrate-down # rollback migrations
```

---

## API Reference

### Auth
| Method | Path | Description | Access |
|--------|------|-------------|--------|
| POST | `/auth/register` | Register a new user | Public |
| POST | `/auth/login` | Login and get JWT | Public |

### Users
| Method | Path | Description | Access |
|--------|------|-------------|--------|
| GET | `/users/me` | Get current user profile | All |
| GET | `/users` | List all users | Admin |
| PATCH | `/users/{id}/block` | Block a user | Admin |
| PATCH | `/users/{id}/unblock` | Unblock a user | Admin |

### Tasks
| Method | Path | Description | Access |
|--------|------|-------------|--------|
| GET | `/tasks` | List tasks | All |
| GET | `/tasks/{id}` | Get task details | All |
| POST | `/tasks` | Create a task | Customer |
| PATCH | `/tasks/{id}/publish` | Publish a task | Customer |
| PATCH | `/tasks/{id}/cancel` | Cancel a task | Customer |
| PATCH | `/tasks/{id}/complete` | Mark task as completed | Executor |
| PATCH | `/tasks/{id}/confirm` | Confirm completion + process payment | Customer |

### Bids
| Method | Path | Description | Access |
|--------|------|-------------|--------|
| GET | `/tasks/{id}/bids` | List bids for a task | All |
| POST | `/tasks/{id}/bids` | Create a bid | Executor |
| PATCH | `/tasks/{id}/bids/{bid_id}/accept` | Accept a bid | Customer |

### Payments
| Method | Path | Description | Access |
|--------|------|-------------|--------|
| GET | `/payments` | Get payment history | All |

---

## Roles

| Role | Description |
|------|-------------|
| `customer` | Creates tasks, selects executors, confirms completion |
| `executor` | Responds to tasks, marks them as completed |
| `admin` | Full access, can block/unblock users |

---

## Task Statuses
```
draft → published → in_progress → completed
  ↓          ↓            ↓
cancelled  cancelled  (cannot cancel)
```

| Status | Description |
|--------|-------------|
| `draft` | Created but not published |
| `published` | Published, accepting bids |
| `in_progress` | Executor selected, work in progress |
| `completed` | Executor marked as done, awaiting confirmation |
| `cancelled` | Cancelled by customer |

---

## curl Examples

### Register a customer
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123","role":"customer"}'
```

### Register an executor
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"executor@example.com","password":"password123","role":"executor"}'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}'
```

### Create a task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"title":"Build a landing page","description":"Need a responsive React landing page","reward":500}'
```

### Publish a task
```bash
curl -X PATCH http://localhost:8080/tasks/<TASK_ID>/publish \
  -H "Authorization: Bearer <TOKEN>"
```

### Create a bid
```bash
curl -X POST http://localhost:8080/tasks/<TASK_ID>/bids \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <EXECUTOR_TOKEN>" \
  -d '{"comment":"I can deliver this in 3 days"}'
```

### Accept a bid
```bash
curl -X PATCH http://localhost:8080/tasks/<TASK_ID>/bids/<BID_ID>/accept \
  -H "Authorization: Bearer <CUSTOMER_TOKEN>"
```

### Mark task as completed
```bash
curl -X PATCH http://localhost:8080/tasks/<TASK_ID>/complete \
  -H "Authorization: Bearer <EXECUTOR_TOKEN>"
```

### Confirm completion
```bash
curl -X PATCH http://localhost:8080/tasks/<TASK_ID>/confirm \
  -H "Authorization: Bearer <CUSTOMER_TOKEN>"
```

### Get my profile
```bash
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer <TOKEN>"
```

### Get payment history
```bash
curl http://localhost:8080/payments \
  -H "Authorization: Bearer <TOKEN>"
```

### Block a user (admin only)
```bash
curl -X PATCH http://localhost:8080/users/<USER_ID>/block \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

---

## Tests
```bash
make test
```

The following business rules are covered:

- Registration with duplicate email
- Login with wrong password
- Login with blocked account
- Task creation success
- Publishing someone else's task → forbidden
- Publishing task with invalid status
- Cancelling a completed task
- Marking completed by wrong executor → forbidden
- Confirming by non-owner → forbidden
- Confirming with insufficient balance
- Bidding on unpublished task
- Duplicate bid → error
- Accepting bid when one already accepted
- Accepting bid as non-owner → forbidden

---

## License

[MIT](LICENSE)