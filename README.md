# Obelisk

A full-stack document management and sharing platform built in Go. Users can upload, download, and share documents with other users through a web dashboard. Features role-based access control with bcrypt-secured authentication.

Built as a hands-on learning project to develop practical Go skills — covering HTTP server design, authentication, database management, and containerised deployment.

## Features

- User registration and login with bcrypt password hashing
- Document upload and download
- Personal document dashboard (`/my-docs`)
- Document sharing between users (`/share`, `/shared-docs`)
- Role-based access control (admin / user)
- Admin user seeding on startup
- CORS middleware for frontend/backend separation
- Dynamic port configuration via environment variable
- Containerised deployment with Docker and Docker Compose

## Tech stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (standard library) |
| Frontend | HTML / JavaScript |
| Auth | bcrypt password hashing |
| Infrastructure | Docker, Docker Compose |

## Getting started

**Prerequisites**

- Docker and Docker Compose

**Run with Docker Compose**

```bash
git clone https://github.com/AlexMcHugh1/go-library-app.git
cd go-library-app
docker-compose up
```

The application will be available at `http://localhost:8080`

**Default admin credentials**

```
Username: admin
Password: admin123
```

> Change the default admin password before deploying to any non-local environment.

**Environment variables**

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |

## API endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/register` | Register a new user |
| `POST` | `/login` | Authenticate and receive token |
| `POST` | `/upload` | Upload a document |
| `GET` | `/list` | List all documents |
| `GET` | `/my-docs` | List documents owned by current user |
| `GET` | `/shared-docs` | List documents shared with current user |
| `GET` | `/download` | Download a document |
| `POST` | `/share` | Share a document with another user |

## Project structure

```
├── cmd/server/           # Application entrypoint (main.go)
├── internal/
│   ├── auth/             # Password hashing and authentication
│   ├── database/         # DB initialisation, migration, queries
│   ├── handlers/         # HTTP request handlers
│   └── models/           # Data models (User, Document etc.)
├── pkg/helpers/          # Shared utility functions
├── assets/diagrams/      # Database ER diagram
├── uploads/              # Uploaded document storage
├── Dockerfile
└── docker-compose.yml
```

## Architecture

See [`assets/diagrams`](assets/diagrams) for the database entity relationship diagram.

The server uses Go's standard `net/http` library with no external web framework. Each handler receives a database connection and manages its own request lifecycle. CORS is handled via middleware wrapping the default ServeMux.

## Why I built this

Built to develop hands-on Go experience beyond tutorials — specifically around HTTP server design using the standard library, authentication flows, database-backed CRUD operations, and containerised deployment. The project evolved from an initial PDF vault concept into a general document sharing platform as scope and understanding grew.

## Security considerations

- Passwords are hashed with bcrypt before storage — plaintext passwords are never persisted
- CORS is currently set to allow all origins (`*`) — restrict to specific origins before any production deployment
- Default admin credentials should be changed before any production deployment

## License

Apache 2.0 — see [LICENSE](LICENSE)
