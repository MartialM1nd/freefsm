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

## Quick Start

1. Clone the repository:
   ```sh
   git clone https://github.com/MartialM1nd/freefsm.git
   cd freefsm
   ```

2. Create the database:
   ```sh
   createdb freefsm
   ```

3. Configure environment:
   ```sh
   cp .env.example .env
   # Edit .env with your PostgreSQL credentials
   ```

4. Run migrations:
   ```sh
   make migrate
   ```

5. Start the server:
   ```sh
   make run
   ```

6. Open http://localhost:8080

## Default Login

- **Email:** admin@freefsm.local
- **Password:** changeme123

**Change this immediately after first login.**

## Installation (Production)

### FreeBSD

```sh
# Build
make build

# Install to /usr/local
make install

# Copy and edit config
cp /usr/local/share/freefsm/freefsm.conf.sample /usr/local/etc/freefsm.conf

# Enable and start service
sysrc freefsm_enable=YES
service freefsm start
```

## Project Structure

```
freefsm/
├── cmd/freefsm/          # Application entry point
├── internal/
│   ├── config/           # Configuration loading
│   ├── database/         # DB connection and migrations
│   ├── handlers/         # HTTP handlers + templates
│   ├── middleware/       # Auth middleware
│   ├── models/           # Data structures
│   └── repository/       # Database queries
├── ui/
│   ├── static/           # CSS, JS (vendored)
│   └── templates/        # HTML templates (if external)
├── deploy/freebsd/       # FreeBSD rc.d scripts
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

## License

MIT
