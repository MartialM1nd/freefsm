package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/models"
)

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if h.cfg.SetupToken == "" {
		http.NotFound(w, r)
		return
	}

	adminExists, err := h.userRepo.AdminExists(r.Context())
	if err != nil || adminExists {
		http.NotFound(w, r)
		return
	}

	h.render(w, r, "pages/register.html", nil)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if h.cfg.SetupToken == "" {
		http.NotFound(w, r)
		return
	}

	adminExists, err := h.userRepo.AdminExists(r.Context())
	if err != nil || adminExists {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.render(w, r, "pages/register.html", map[string]any{
			"Error": "Invalid form data",
		})
		return
	}

	token := r.FormValue("token")
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	if token != h.cfg.SetupToken {
		h.render(w, r, "pages/register.html", map[string]any{
			"Error": "Invalid setup token",
			"Email": email,
			"Name":  name,
			"Phone": phone,
		})
		return
	}

	if email == "" || password == "" || name == "" {
		h.render(w, r, "pages/register.html", map[string]any{
			"Error": "Email, password, and name are required",
			"Email": email,
			"Name":  name,
			"Phone": phone,
		})
		return
	}

	user := &models.User{
		Email: email,
		Name:  name,
		Phone: phone,
		Role:  "admin",
	}

	if err := h.userRepo.Create(r.Context(), user, password); err != nil {
		h.render(w, r, "pages/register.html", map[string]any{
			"Error": "Failed to create user: " + err.Error(),
			"Email": email,
			"Name":  name,
			"Phone": phone,
		})
		return
	}

	session, err := h.userRepo.CreateSession(r.Context(), user.ID)
	if err != nil {
		h.render(w, r, "pages/register.html", map[string]any{
			"Error": "Failed to create session",
			"Email": email,
			"Name":  name,
			"Phone": phone,
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}