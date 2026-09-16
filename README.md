# Go Note API

A RESTful Notes API built with Go and MySQL, featuring JWT-based
authentication, layered architecture, and a React/TypeScript frontend.

## Features

- CRUD operations for notes
- JWT authentication with bcrypt password hashing
- Layered architecture (handler → service → repository)
- Pagination support
- Centralized error handling
- Unit tests
- Dockerized development environment

## Tech Stack

**Backend:** Go, MySQL, JWT, bcrypt, Docker — hand-written
**Frontend:** React, TypeScript, Vite, Tailwind CSS — built with an AI
coding agent as a lightweight client for the API

## Getting Started

### Backend

```bash
go run cmd/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## API Endpoints

| Method | Endpoint        | Description         |
|--------|-----------------|----------------------|
| POST   | `/register`     | Register a new user  |
| POST   | `/login`        | Log in               |
| POST   | `/notes`        | Create a note        |
| GET    | `/notes`        | List notes           |
| GET    | `/notes/get`    | Get a note by ID     |
| PUT    | `/notes/update` | Update a note        |
| DELETE | `/notes/delete` | Delete a note        |

## License

MIT
