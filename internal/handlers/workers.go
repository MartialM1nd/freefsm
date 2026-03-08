package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) WorkersList(w http.ResponseWriter, r *http.Request) {
	workers, err := h.userRepo.List(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load workers")
		return
	}

	h.render(w, r, "pages/workers/list.html", map[string]any{
		"Title":   "Workers",
		"Workers": workers,
	})
}

func (h *Handler) WorkersNew(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "pages/workers/form.html", map[string]any{
		"Title":  "New Worker",
		"Worker": &models.User{Role: "technician"},
		"IsNew":  true,
	})
}

func (h *Handler) WorkersCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	worker := &models.User{
		Email: r.FormValue("email"),
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
		Role:  r.FormValue("role"),
	}
	password := r.FormValue("password")

	if worker.Name == "" || worker.Email == "" {
		h.render(w, r, "pages/workers/form.html", map[string]any{
			"Title":  "New Worker",
			"Worker": worker,
			"IsNew":  true,
			"Error":  "Name and email are required",
		})
		return
	}

	if password == "" {
		h.render(w, r, "pages/workers/form.html", map[string]any{
			"Title":  "New Worker",
			"Worker": worker,
			"IsNew":  true,
			"Error":  "Password is required",
		})
		return
	}

	if worker.Role == "" {
		worker.Role = "technician"
	}

	if err := h.userRepo.Create(r.Context(), worker, password); err != nil {
		h.render(w, r, "pages/workers/form.html", map[string]any{
			"Title":  "New Worker",
			"Worker": worker,
			"IsNew":  true,
			"Error":  "Failed to create worker. Email may already be in use.",
		})
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/workers")
		return
	}
	http.Redirect(w, r, "/workers", http.StatusSeeOther)
}

func (h *Handler) WorkersView(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	worker, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load worker")
		return
	}
	if worker == nil {
		h.errorResponse(w, http.StatusNotFound, "Worker not found")
		return
	}

	// Get assigned jobs
	jobs, _ := h.jobRepo.ListByWorker(r.Context(), id)

	h.render(w, r, "pages/workers/view.html", map[string]any{
		"Title":  worker.Name,
		"Worker": worker,
		"Jobs":   jobs,
	})
}

func (h *Handler) WorkersEdit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	worker, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load worker")
		return
	}
	if worker == nil {
		h.errorResponse(w, http.StatusNotFound, "Worker not found")
		return
	}

	h.render(w, r, "pages/workers/form.html", map[string]any{
		"Title":  "Edit " + worker.Name,
		"Worker": worker,
		"IsNew":  false,
	})
}

func (h *Handler) WorkersUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	worker := &models.User{
		ID:    id,
		Email: r.FormValue("email"),
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
		Role:  r.FormValue("role"),
	}

	if worker.Name == "" || worker.Email == "" {
		h.render(w, r, "pages/workers/form.html", map[string]any{
			"Title":  "Edit Worker",
			"Worker": worker,
			"IsNew":  false,
			"Error":  "Name and email are required",
		})
		return
	}

	if err := h.userRepo.Update(r.Context(), worker); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to update worker")
		return
	}

	// Update password if provided
	if password := r.FormValue("password"); password != "" {
		if err := h.userRepo.UpdatePassword(r.Context(), id, password); err != nil {
			h.errorResponse(w, http.StatusInternalServerError, "Failed to update password")
			return
		}
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/workers/"+id.String())
		return
	}
	http.Redirect(w, r, "/workers/"+id.String(), http.StatusSeeOther)
}

func (h *Handler) WorkersDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	if err := h.userRepo.Delete(r.Context(), id); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to delete worker")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/workers")
		return
	}
	http.Redirect(w, r, "/workers", http.StatusSeeOther)
}
