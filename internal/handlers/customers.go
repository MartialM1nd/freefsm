package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CustomersList(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customerRepo.List(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load customers")
		return
	}

	h.render(w, r, "pages/customers/list.html", map[string]any{
		"Title":     "Customers",
		"Customers": customers,
	})
}

func (h *Handler) CustomersNew(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "pages/customers/form.html", map[string]any{
		"Title":    "New Customer",
		"Customer": &models.Customer{},
		"IsNew":    true,
	})
}

func (h *Handler) CustomersCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	customer := &models.Customer{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
		Address: r.FormValue("address"),
		City:    r.FormValue("city"),
		State:   r.FormValue("state"),
		Zip:     r.FormValue("zip"),
		Notes:   r.FormValue("notes"),
	}

	if customer.Name == "" {
		h.render(w, r, "pages/customers/form.html", map[string]any{
			"Title":    "New Customer",
			"Customer": customer,
			"IsNew":    true,
			"Error":    "Name is required",
		})
		return
	}

	if err := h.customerRepo.Create(r.Context(), customer); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/customers")
		return
	}
	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

func (h *Handler) CustomersView(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.customerRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load customer")
		return
	}
	if customer == nil {
		h.errorResponse(w, http.StatusNotFound, "Customer not found")
		return
	}

	h.render(w, r, "pages/customers/view.html", map[string]any{
		"Title":    customer.Name,
		"Customer": customer,
	})
}

func (h *Handler) CustomersEdit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.customerRepo.GetByID(r.Context(), id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to load customer")
		return
	}
	if customer == nil {
		h.errorResponse(w, http.StatusNotFound, "Customer not found")
		return
	}

	h.render(w, r, "pages/customers/form.html", map[string]any{
		"Title":    "Edit " + customer.Name,
		"Customer": customer,
		"IsNew":    false,
	})
}

func (h *Handler) CustomersUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	customer := &models.Customer{
		ID:      id,
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
		Address: r.FormValue("address"),
		City:    r.FormValue("city"),
		State:   r.FormValue("state"),
		Zip:     r.FormValue("zip"),
		Notes:   r.FormValue("notes"),
	}

	if customer.Name == "" {
		h.render(w, r, "pages/customers/form.html", map[string]any{
			"Title":    "Edit Customer",
			"Customer": customer,
			"IsNew":    false,
			"Error":    "Name is required",
		})
		return
	}

	if err := h.customerRepo.Update(r.Context(), customer); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to update customer")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/customers/"+id.String())
		return
	}
	http.Redirect(w, r, "/customers/"+id.String(), http.StatusSeeOther)
}

func (h *Handler) CustomersDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := h.customerRepo.Delete(r.Context(), id); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, "Failed to delete customer")
		return
	}

	if h.isHTMX(r) {
		w.Header().Set("HX-Redirect", "/customers")
		return
	}
	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}
