package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get today's jobs
	today := time.Now().Truncate(24 * time.Hour)
	jobs, err := h.jobRepo.ListByDate(r.Context(), today)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load jobs")
		return
	}

	// Get counts by status
	allJobs, _ := h.jobRepo.List(r.Context())
	statusCounts := make(map[string]int)
	for _, j := range allJobs {
		statusCounts[string(j.Status)]++
	}

	h.render(w, r, "pages/dashboard.html", map[string]any{
		"Title":        "Dashboard",
		"Jobs":         jobs,
		"Today":        today,
		"StatusCounts": statusCounts,
	})
}
