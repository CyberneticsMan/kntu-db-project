# KNTU DB Project

Go + Gin backend for a sports ticket reservation platform.

## Features implemented

- User registration and login with JWT
- User CRUD endpoints
- Ticket CRUD endpoints
- Ticket listing and simple search/filtering
- Reservation creation, cancelation, and payment endpoints
- Ticket issue reporting endpoints
- PostgreSQL migrations for users, tickets, reservations, payments, and reports
- New users created through registration or `/users` default to the `customer` role

## Environment

Copy `.env` and adjust the values if needed:

- `DATABASE_URL`
- `JWT_SECRET`

## Database setup

Run the migrations with the migration command in `cmd/migrate`.

## API routes

Base path: `/api/v1`

### Auth

- `POST /auth/register`
- `POST /auth/login`

### Users

- `GET /users`
- `GET /users/:id`
- `POST /users`
- `PUT /users/:id`
- `DELETE /users/:id`

### Tickets

- `GET /tickets`
- `GET /tickets/:id`
- `POST /tickets`
- `PUT /tickets/:id`
- `DELETE /tickets/:id`

### Reservations

- `POST /reservations`
- `GET /reservations/:id`
- `GET /users/:userId/reservations`
- `POST /reservations/:id/cancel`
- `POST /reservations/:id/pay`

### Reports

- `POST /reports`
- `GET /reports/:id`
- `GET /users/:userId/reports`

## Development

- Build: `go build ./...`
- Test: `go test ./...`