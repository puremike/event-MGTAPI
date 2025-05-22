# Event Management API

A RESTful API for managing events, users, and attendees, built with Go, Gin, and Redis.

## Features

- **User Authentication**: Register, login, and retrieve user details.
- **Event Management**: Create, update, delete, and list events.
- **Attendee Management**: Add or remove attendees from events, and list events for a user.
- **JWT-based Auth**: Secure endpoints with JWT authentication.
- **Redis Caching**: Improve performance for event and user data.
- **Swagger Docs**: API documentation available via Swagger UI.
- **Health & Debug Endpoints**: For monitoring and debugging.

## API Endpoints

### Auth

- `POST /api/v1/auth/register` — Register a new user
- `POST /api/v1/auth/login` — Login and receive JWT
- `GET /api/v1/auth/:id` — Get user by ID

### Events

- `GET /api/v1/events/` — List all events
- `GET /api/v1/events/:id` — Get event by ID
- `GET /api/v1/events/:id/attendees` — List attendees for an event
- `POST /api/v1/events` — Create event (auth required)
- `PUT /api/v1/events/:id` — Update event (auth + event context)
- `DELETE /api/v1/events/:id` — Delete event (auth + event context)

### Attendees

- `GET /api/v1/attendees/:userId/events` — List events a user is attending
- `POST /api/v1/events/:id/attendees/:userId` — Add attendee to event (auth + event context)
- `DELETE /api/v1/events/:id/attendees/:userId` — Remove attendee from event (auth + event context)

### Monitoring

- `GET /api/v1/health` — Health check
- `GET /api/v1/debug/vars` — Debug variables

### Swagger

- `GET /swagger/index.html` — Swagger UI

## Getting Started

### Prerequisites

- Go 1.19+
- Docker (optional, for running with Docker Compose)
- PostgreSQL (or your configured DB)
- Redis

### Setup

1. **Clone the repo:**
   ```bash
   git clone <your-repo-url>
   cd event-MGTAPI
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables:**
   - Copy `.env.example` to `.env` and set your DB, Redis credentials, and JWT secret.

4. **Run database migrations:**
   ```bash
   make migrate-up
   ```

5. **Start the server:**
   ```bash
   go run ./cmd/api
   ```

6. **Access Swagger docs:**
   - Visit `http://localhost:8080/swagger/index.html`

### Using Docker

```bash
docker-compose up --build
```

## Project Structure

- `cmd/api/` — Main API application (routes, handlers, middlewares)
- `internal/` — Internal packages (auth, db, storage, cache)
- `docs/` — Swagger/OpenAPI docs
- `migrate/` — Database migration files

## Redis Usage

- Redis is used for caching event and user data to improve performance.
- Configuration is handled via environment variables (see `.env.example`).

## Contributing

Pull requests are welcome! Please open issues for suggestions or bugs.

---

Let me know if you want to add usage examples, environment variables, or anything else!
