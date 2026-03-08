# FreeFSM Agent Guidelines

## Build, Test, and Development Commands

### Makefile Targets

```bash
make build          # Compile binary to ./freefsm
make run            # Run (requires DATABASE_URL, SESSION_SECRET)
make migrate        # Run database migrations
make clean          # Remove compiled binary
make install        # Install to PREFIX (/usr/local by default)
make uninstall      # Remove installed files
```

### Go Commands

```bash
go build -o freefsm ./cmd/freefsm          # Build binary
go run ./cmd/freefsm -migrate               # Run migrations
go test ./...                               # Run all tests
go test -v -run TestFunctionName ./...      # Run single test by name
go test -v -run "^TestUser" ./...           # Run tests matching regex
go test -count=1 ./...                      # Run tests without cache
```

### Environment Variables

- Required: `DATABASE_URL`, `SESSION_SECRET`
- Optional: `PORT` (default: 8080), `ENV` (default: development), `STATIC_PATH` (default: ui/static)

## Code Style Guidelines

### Import Ordering

Three groups with blank lines: stdlib, third-party, internal packages.

```go
import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/MartialM1nd/freefsm/internal/config"
	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/MartialM1nd/freefsm/internal/repository"
)
```

### Naming

- Structs: PascalCase (`User`, `JobRepo`, `Handler`)
- Methods: PascalCase (`GetByID`, `IsAdmin`)
- Variables/Fields: PascalCase exported, camelCase unexported
- Packages: lowercase, single word

### Repository Pattern

```go
type UserRepo struct { db *database.DB }
func NewUserRepo(db *database.DB) *UserRepo { return &UserRepo{db: db} }
func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) { /*...*/ }
```

- All methods take `context.Context` as first param
- Return pointer to model (or nil) plus error
- Use raw SQL with pgx (no ORM)

### Error Handling

```go
err := r.db.Pool.QueryRow(ctx, query, args...).Scan(...)
if err != nil {
	if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
	return nil, err
}
```

- Use `errors.Is()` for comparison
- Return errors directly; wrap only when adding context

### Struct Tags

```go
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`  // Never serialize
	Phone        string     `json:"phone,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
```

### Database Queries

```go
rows, err := r.db.Pool.Query(ctx, `SELECT id, email, name FROM users WHERE deleted_at IS NULL`)
if err != nil { return nil, err }
defer rows.Close()
```
- Use uppercase SQL keywords, `$1, $2, ...` placeholders
- Always defer rows.Close()

### Templates

Embed with `//go:embed` in `internal/handlers/templates/`:

```go
//go:embed templates
var templatesFS embed.FS
```

Define template functions in `loadTemplates()`:

```go
funcMap := template.FuncMap{
	"statusClass": func(status models.JobStatus) string { ... },
}
```

### Handlers

Use `render()` for HTMX + full page support:

```go
func (h *Handler) render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) {
	if r.Header.Get("HX-Request") == "true" {
		h.templates.ExecuteTemplate(w, name, data)
	} else {
		data["Content"] = name
		h.templates.ExecuteTemplate(w, "layouts/base.html", data)
	}
}
```

### Middleware

- Auth middleware in `internal/middleware/auth.go`
- Store user in context via `context.WithValue`
- Retrieve with `middleware.GetUser(r.Context())`

### Configuration

- Use `config.Config` struct
- Load from env vars with defaults
- Validate required fields on load

## Project Structure

```
internal/
├── config/       # Configuration
├── database/     # DB connection + migrations
├── handlers/     # HTTP handlers + templates
├── middleware/   # Auth
├── models/       # Data structures
└── repository/   # DB queries

cmd/freefsm/      # Entry point
ui/static/        # Vendored CSS/JS
deploy/freebsd/   # rc.d script + config
```

## FreeBSD Notes

- Use dedicated `freefsm` user (not www or nobody)
- rc.d to `/usr/local/etc/rc.d/`
- Config in `/usr/local/share/freefsm/`