package handlers

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/MartialM1nd/freefsm/internal/config"
	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/middleware"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/MartialM1nd/freefsm/internal/repository"
)

//go:embed templates
var templatesFS embed.FS

type Handler struct {
	db           *database.DB
	cfg          *config.Config
	templates    *template.Template
	userRepo     *repository.UserRepo
	customerRepo *repository.CustomerRepo
	jobRepo      *repository.JobRepo
}

func New(db *database.DB, cfg *config.Config) *Handler {
	h := &Handler{
		db:           db,
		cfg:          cfg,
		userRepo:     repository.NewUserRepo(db),
		customerRepo: repository.NewCustomerRepo(db),
		jobRepo:      repository.NewJobRepo(db),
	}

	h.loadTemplates()
	return h
}

func (h *Handler) loadTemplates() {
	funcMap := template.FuncMap{
		"statusClass": func(status models.JobStatus) string {
			classes := map[models.JobStatus]string{
				models.JobStatusNew:             "secondary",
				models.JobStatusInTransit:       "primary",
				models.JobStatusInProgress:      "primary",
				models.JobStatusPending:         "warning",
				models.JobStatusScheduledReturn: "warning",
				models.JobStatusReadyToInvoice:  "success",
				models.JobStatusCompleted:       "success",
				models.JobStatusCancelled:       "error",
			}
			return classes[status]
		},
		"priorityClass": func(priority models.JobPriority) string {
			classes := map[models.JobPriority]string{
				models.JobPriorityLow:    "secondary",
				models.JobPriorityMedium: "primary",
				models.JobPriorityHigh:   "warning",
				models.JobPriorityUrgent: "error",
			}
			return classes[priority]
		},
	}

	h.templates = template.New("").Funcs(funcMap)

	// Walk embedded templates and register with stripped prefix
	fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		// Strip "templates/" prefix for cleaner template names
		name := strings.TrimPrefix(path, "templates/")

		content, err := templatesFS.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read template %s: %v", path, err)
		}

		_, err = h.templates.New(name).Parse(string(content))
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}

		return nil
	})
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) {
	if data == nil {
		data = make(map[string]any)
	}

	data["User"] = middleware.GetUser(r.Context())

	if r.Header.Get("HX-Request") == "true" {
		if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
			log.Printf("Template error (partial %s): %v", name, err)
			h.errorResponse(w, 500, "Internal server error")
		}
	} else {
		var buf strings.Builder
		if err := h.templates.ExecuteTemplate(&buf, name, data); err != nil {
			log.Printf("Template error (content %s): %v", name, err)
			h.errorResponse(w, 500, "Internal server error")
			return
		}
		data["Content"] = template.HTML(buf.String())
		if err := h.templates.ExecuteTemplate(w, "layouts/base.html", data); err != nil {
			log.Printf("Template error (layout): %v", err)
			h.errorResponse(w, 500, "Internal server error")
		}
	}
}

func (h *Handler) renderPartial(w http.ResponseWriter, name string, data any) {
	h.templates.ExecuteTemplate(w, name, data)
}

func (h *Handler) renderTemplate(w io.Writer, name string, data any) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

func (h *Handler) parseTemplates(pattern string) *template.Template {
	return template.Must(template.ParseGlob(filepath.Join("ui/templates", pattern)))
}

func (h *Handler) errorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func (h *Handler) isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
