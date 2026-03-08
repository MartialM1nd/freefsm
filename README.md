# FreeFSM

A simple, high-performance Field Service Management system built with FreeBSD in mind.

## Features

- Job/Work Order Management
- Customer Management
- Technician/Worker Management
- Job Scheduling with Calendar View
- Job Notes and History Tracking
- Status Workflow: New → In Transit → In Progress → Pending → Scheduled Return → Ready to Invoice → Completed

## Tech Stack

- **Backend:** Go 1.21+ with Chi router
- **Database:** PostgreSQL 15+
- **Frontend:** HTMX + Alpine.js
- **CSS:** Pico CSS (no build step)
- **Templating:** Go html/template

## Prerequisites

- Go 1.21+
- PostgreSQL 15+

### FreeBSD

```sh
pkg install go postgresql15-server
```

### Linux (Fedora/RHEL)

```sh
dnf install golang postgresql-server
```

## Quick Start (Development)

1. Clone the repository:
   ```sh
   git clone https://github.com/MartialM1nd/freefsm.git
   cd freefsm
   ```

2. Set up PostgreSQL database and user:
   ```sh
   # Create database user (will prompt for password)
   createuser -P freefsm

   # Create database owned by that user
   createdb -O freefsm freefsm
   ```

3. Configure environment:
   ```sh
   cp .env.example .env
   ```

   Edit `.env` and update the database connection string with your password:
   ```
   DATABASE_URL=postgres://freefsm:yourpassword@localhost:5432/freefsm?sslmode=disable
   ```

   Generate a session secret:
   ```sh
   openssl rand -hex 32
   ```

4. Get dependencies:
   ```sh
   go mod tidy
   ```

5. Run migrations:
   ```sh
   make migrate
   ```

6. Start the server:
   ```sh
   make run
   ```

7. Open http://localhost:8080

## Default Login

- **Email:** admin@freefsm.local
- **Password:** changeme123

**Change this immediately after first login.**

## Installation (Production - FreeBSD)

1. Create service user:
   ```sh
   pw useradd -n freefsm -c "FreeFSM Service" -s /usr/sbin/nologin -w no
   ```

2. Set up PostgreSQL database:
   ```sh
   # As postgres user or with appropriate privileges
   createuser -P freefsm
   createdb -O freefsm freefsm
   ```

3. Build and install:
   ```sh
   make build
   make install
   ```

4. Configure:
   ```sh
   cp /usr/local/share/freefsm/freefsm.conf.sample /usr/local/etc/freefsm.conf
   ```

   Edit `/usr/local/etc/freefsm.conf`:
   - Set `DATABASE_URL` with your PostgreSQL credentials
   - Generate and set `SESSION_SECRET` with `openssl rand -hex 32`

5. Run migrations:
   ```sh
   # Set environment temporarily
   . /usr/local/etc/freefsm.conf
   export DATABASE_URL SESSION_SECRET
   /usr/local/bin/freefsm -migrate
   ```

6. Enable and start service:
   ```sh
   sysrc freefsm_enable=YES
   service freefsm start
   ```

7. Verify it's running:
   ```sh
   service freefsm status
   ```

## Project Structure

```
freefsm/
├── cmd/freefsm/          # Application entry point
├── internal/
│   ├── config/           # Configuration loading
│   ├── database/         # DB connection and migrations
│   ├── handlers/         # HTTP handlers + embedded templates
│   ├── middleware/       # Auth middleware
│   ├── models/           # Data structures
│   └── repository/       # Database queries
├── ui/static/            # CSS, JS (vendored)
├── deploy/freebsd/       # FreeBSD rc.d script and config
├── Makefile
└── .env.example
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | (required) |
| `SESSION_SECRET` | Secret key for session cookies | (required) |
| `PORT` | HTTP server port | 8080 |
| `ENV` | Environment (development/production) | development |
| `STATIC_PATH` | Path to static files | ui/static |

## License

MIT
