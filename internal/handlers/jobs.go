package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MartialM1nd/freefsm/internal/middleware"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) JobsList(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobRepo.List(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load jobs")
		return
	}

	h.render(w, r, "pages/jobs/list.html", map[string]any{
		"Title": "Jobs",
		"Jobs":  jobs,
	})
}

func (h *Handler) JobsNew(w http.ResponseWriter, r *http.Request) {
	customers, _ := h.customerRepo.List(r.Context())
	workers, _ := h.userRepo.ListTechnicians(r.Context())

	h.render(w, r, "pages/jobs/form.html", map[string]any{
		"Title":     "New Job",
		"Job":       &models.Job{Status: models.JobStatusNew, Priority: models.JobPriorityMedium, UseCustomerAddress: true},
		"IsNew":     true,
		"Customers": customers,
		"Workers":   workers,
		"Statuses":  jobStatuses(),
		"Priorities": jobPriorities(),
	})
}

func (h *Handler) JobsCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	job, err := h.parseJobForm(r)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if job.Title == "" {
		customers, _ := h.customerRepo.List(r.Context())
		workers, _ := h.userRepo.ListTechnicians(r.Context())
		h.render(w, r, "pages/jobs/form.html", map[string]any{
			"Title":      "New Job",
			"Job":        job,
			"IsNew":      true,
			"Customers":  customers,
			"Workers":    workers,
			"Statuses":   jobStatuses(),
			"Priorities": jobPriorities(),
			"Error":      "Title is required",
		})
		return
	}

	if err := h.jobRepo.Create(r.Context(), job); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to create job")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/jobs")
		return
	}
	http.Redirect(w, r, "/jobs", http.StatusSeeOther)
}

func (h *Handler) JobsView(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.jobRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load job")
		return
	}
	if job == nil {
		h.errorResponse(w, http.StatusNotFound, "Job not found")
		return
	}

	// Load related data
	if job.CustomerID != nil {
		job.Customer, _ = h.customerRepo.GetByID(r.Context(), *job.CustomerID)
	}
	if job.AssignedTo != nil {
		job.AssignedUser, _ = h.userRepo.GetByID(r.Context(), *job.AssignedTo)
	}
	notes, _ := h.jobRepo.GetNotes(r.Context(), id)
	history, _ := h.jobRepo.GetHistory(r.Context(), id)

	h.render(w, r, "pages/jobs/view.html", map[string]any{
		"Title":    job.Title,
		"Job":      job,
		"Notes":    notes,
		"History":  history,
		"Statuses": jobStatuses(),
	})
}

func (h *Handler) JobsEdit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.jobRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load job")
		return
	}
	if job == nil {
		h.errorResponse(w, http.StatusNotFound, "Job not found")
		return
	}

	customers, _ := h.customerRepo.List(r.Context())
	workers, _ := h.userRepo.ListTechnicians(r.Context())

	h.render(w, r, "pages/jobs/form.html", map[string]any{
		"Title":      "Edit " + job.Title,
		"Job":        job,
		"IsNew":      false,
		"Customers":  customers,
		"Workers":    workers,
		"Statuses":   jobStatuses(),
		"Priorities": jobPriorities(),
	})
}

func (h *Handler) JobsUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	job, err := h.parseJobForm(r)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	job.ID = id

	if job.Title == "" {
		customers, _ := h.customerRepo.List(r.Context())
		workers, _ := h.userRepo.ListTechnicians(r.Context())
		h.render(w, r, "pages/jobs/form.html", map[string]any{
			"Title":      "Edit Job",
			"Job":        job,
			"IsNew":      false,
			"Customers":  customers,
			"Workers":    workers,
			"Statuses":   jobStatuses(),
			"Priorities": jobPriorities(),
			"Error":      "Title is required",
		})
		return
	}

	if err := h.jobRepo.Update(r.Context(), job); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to update job")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/jobs/"+id.String())
		return
	}
	http.Redirect(w, r, "/jobs/"+id.String(), http.StatusSeeOther)
}

func (h *Handler) JobsDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	if err := h.jobRepo.Delete(r.Context(), id); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to delete job")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/jobs")
		return
	}
	http.Redirect(w, r, "/jobs", http.StatusSeeOther)
}

func (h *Handler) JobsAddNote(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	user := middleware.GetUser(r.Context())
	content := r.FormValue("content")
	if content == "" {
		h.errorResponse(w, http.StatusBadRequest, "Note content is required")
		return
	}

	note := &models.JobNote{
		JobID:   id,
		UserID:  user.ID,
		Content: content,
	}

	if err := h.jobRepo.AddNote(r.Context(), note); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to add note")
		return
	}

	// Return updated notes list for HTMX
	if h.isHTMX(r) {
		notes, _ := h.jobRepo.GetNotes(r.Context(), id)
		h.renderPartial(w, "partials/job_notes.html", map[string]any{
			"Notes": notes,
			"JobID": id,
		})
		return
	}

	http.Redirect(w, r, "/jobs/"+id.String(), http.StatusSeeOther)
}

func (h *Handler) JobsUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	user := middleware.GetUser(r.Context())
	status := models.JobStatus(r.FormValue("status"))

	if err := h.jobRepo.UpdateStatus(r.Context(), id, status, user.ID); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	if h.isHTMX(r) {
		job, _ := h.jobRepo.GetByID(r.Context(), id)
		h.renderPartial(w, "partials/job_status.html", map[string]any{
			"Job":      job,
			"Statuses": jobStatuses(),
		})
		return
	}

	http.Redirect(w, r, "/jobs/"+id.String(), http.StatusSeeOther)
}

func (h *Handler) parseJobForm(r *http.Request) (*models.Job, error) {
	job := &models.Job{
		Title:              r.FormValue("title"),
		Description:        r.FormValue("description"),
		Status:             models.JobStatus(r.FormValue("status")),
		Priority:           models.JobPriority(r.FormValue("priority")),
		UseCustomerAddress: r.FormValue("use_customer_address") == "on",
		LocationAddress:    r.FormValue("location_address"),
		LocationCity:       r.FormValue("location_city"),
		LocationState:      r.FormValue("location_state"),
		LocationZip:        r.FormValue("location_zip"),
		ScheduledTime:      r.FormValue("scheduled_time"),
	}

	// Customer
	if customerID := r.FormValue("customer_id"); customerID != "" {
		id, err := uuid.Parse(customerID)
		if err == nil {
			job.CustomerID = &id
		}
	}

	// Assigned worker
	if assignedTo := r.FormValue("assigned_to"); assignedTo != "" {
		id, err := uuid.Parse(assignedTo)
		if err == nil {
			job.AssignedTo = &id
		}
	}

	// Scheduled date
	if dateStr := r.FormValue("scheduled_date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			job.ScheduledDate = &date
		}
	}

	// Estimated duration
	if durationStr := r.FormValue("estimated_duration"); durationStr != "" {
		if duration, err := strconv.Atoi(durationStr); err == nil {
			job.EstimatedDuration = &duration
		}
	}

	// Defaults
	if job.Status == "" {
		job.Status = models.JobStatusNew
	}
	if job.Priority == "" {
		job.Priority = models.JobPriorityMedium
	}

	return job, nil
}

func jobStatuses() []struct {
	Value string
	Label string
} {
	return []struct {
		Value string
		Label string
	}{
		{"new", "New"},
		{"in_transit", "In Transit"},
		{"in_progress", "In Progress"},
		{"pending", "Pending"},
		{"scheduled_return", "Scheduled Return"},
		{"ready_to_invoice", "Ready to Invoice"},
		{"completed", "Completed"},
		{"cancelled", "Cancelled"},
	}
}

func jobPriorities() []struct {
	Value string
	Label string
} {
	return []struct {
		Value string
		Label string
	}{
		{"low", "Low"},
		{"medium", "Medium"},
		{"high", "High"},
		{"urgent", "Urgent"},
	}
}
