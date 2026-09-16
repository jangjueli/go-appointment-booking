# Hospital Appointment Booking Module

**Author:** Dacha Jirawutthiwongchai

Backend service for managing doctor schedules, appointment types, appointment bookings, and available appointment slots for a hospital outpatient department.

## Tech Stack

- Go 1.26+
- Gin
- GORM
- PostgreSQL
- Docker / Docker Compose
- Go testing package

## Features

- Doctor schedule management
- Appointment type management
- Appointment booking
- Available slot calculation
- Break time handling
- Overlapping appointment validation
- PostgreSQL persistence
- Database seeding
- Unit tests

## Project Structure

```text
go-appointment-booking/
├── cmd/api/
│   └── main.go
├── internal/
│   ├── database/
│   ├── handler/
│   ├── model/
│   └── service/
├── docker-compose.yml
├── .env
├── go.mod
├── go.sum
└── README.md
```

## Requirements

- Docker
- Docker Compose
Go 1.26+ is only required if you want to run the application locally without Docker.

## Configuration

For local development, create `.env` in the project root:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointment

When running with Docker Compose, the API connects to PostgreSQL using the Docker service name:

DB_HOST=postgres
DB_PORT=5432

## Run with Docker Compose

Build and start all services:

```bash
docker compose up -d --build
```

Services:

API      http://localhost:8080
pgAdmin  http://localhost:8081
Postgres localhost:5432

The application automatically runs database migration and seed data on startup.

View API logs:

```bash
docker compose logs -f api
```

Stop services:

```bash
docker compose down
```

To remove the database volume and start with a clean database:

```bash
docker compose down -v
docker compose up -d --build
```

pgAdmin:

```text
http://localhost:8081
```

Login:

Email:    root@gmail.com
Password: 1234

When adding the PostgreSQL server in pgAdmin, use:

Host:     postgres
Port:     5432
Database: appointment
Username: postgres
Password: postgres

## API Endpoints

### Doctor Schedule

```http
POST /doctors/:doctorID/schedules
GET  /doctors/:doctorID/schedules
```

Example request:

```json
{
  "day_of_week": 1,
  "start_time": "09:00",
  "end_time": "17:00",
  "break_start_time": "12:00",
  "break_end_time": "13:00",
  "accept_booking": true
}
```

`day_of_week`:

```text
0 = Sunday
1 = Monday
2 = Tuesday
3 = Wednesday
4 = Thursday
5 = Friday
6 = Saturday
```

### Appointment Types

```http
GET /appointment-types
```

Default appointment types:

Type                           Duration
--------------------------------------
New patient visit               60 min
Follow-up visit                 30 min
Consultation                    30 min
Procedure                       60 min

### Book Appointment

```http
POST /appointments
```

Example:

```json
{
  "patient_id": 1,
  "doctor_id": 1,
  "department_id": 1,
  "appointment_type_id": 2,
  "appointment_date": "2026-09-21",
  "start_time": "10:30",
  "reason": "Follow-up visit",
  "created_by": "staff01"
}
```

The appointment `end_time` is calculated automatically from the appointment type duration.

### Available Slots

```http
GET /doctors/:doctorID/available-slots
```

Query parameters:

```text
date
appointment_type_id
```

Example:

```text
GET /doctors/1/available-slots?date=2026-09-21&appointment_type_id=2
```

Available slots are calculated from:

- Doctor working hours
- Appointment type duration
- Break time
- Existing booked appointments
- Accept booking

## Validation

When creating an appointment, the service validates:

- Patient exists
- Doctor exists
- Department exists
- Appointment type exists
- Doctor belongs to the specified department
- Doctor has a bookable schedule
- Appointment is within working hours
- Appointment does not overlap the break
- Appointment does not overlap an existing booked appointment

## Time Handling

Schedule and appointment times use `HH:mm` strings such as:

```text
09:00
12:30
17:00
```

They are parsed into Go time values only when performing time calculations.

This avoids unnecessary date and timezone handling for time-of-day fields.

## Testing

Unit tests cover:

- Slot generation
- Break time filtering
- Time overlap detection

Run tests:

```bash
go test ./...
```

## Assumptions

- Each doctor has one schedule per day of the week.
- Appointment duration is defined by `AppointmentType`.
- Only appointments with status `booked` block available slots.
- A schedule with `accept_booking = false` is not available for booking.
- Adjacent appointments are allowed, for example `10:00–10:30` followed by `10:30–11:00`.
- The implementation focuses on the requested booking and availability functionality.

## Known Limitations

- Authentication and authorization are not implemented because they are outside the requested scope.
- Appointment cancellation and completion are not implemented because the implementation focuses on requirements 1–4.
- The system assumes a single schedule per doctor per day.
- Appointment start times are accepted based on working-hour and overlap validation; configurable slot intervals are not implemented.
- Database migration and seed data run automatically when the application starts.
- Concurrent booking protection at the database transaction/locking level is not implemented.
