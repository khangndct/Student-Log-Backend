# Student Diary Backend API

A REST API backend for a student diary application that allows users to create, manage, and share diary entries with role-based access control.

## Tech Stack

- **Backend**: Go with Echo framework
- **Database**: PostgreSQL with GORM ORM

## How to Run Locally

1. Clone the repository
2. Create a `.env` file in the root directory with the required environment variables (see below)
3. Run the application using Docker Compose:

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

## Environment Variables

- `POSTGRES_USER` - PostgreSQL username
- `POSTGRES_PASSWORD` - PostgreSQL password
- `POSTGRES_HOST` - PostgreSQL host (use `db` when running with Docker Compose)
- `JWT_SECRET` - Secret key for JWT token signing

Optional:

- `POSTGRES_PORT` - PostgreSQL port (default: 5432)
- `POSTGRES_DB` - PostgreSQL database name (default: postgres)
- `POSTGRES_SSLMODE` - PostgreSQL SSL mode (default: disable)
- `POSTGRES_TIMEZONE` - PostgreSQL timezone (default: Asia/Ho_Chi_Minh)
