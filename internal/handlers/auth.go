package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	if cookie, err := r.Cookie("session"); err == nil {
		if session, _ := h.userRepo.GetSession(r.Context(), cookie.Value); session != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	h.templates.ExecuteTemplate(w, "pages/login.html", nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.templates.ExecuteTemplate(w, "pages/login.html", map[string]any{
			"Error": "Invalid form data",
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		h.templates.ExecuteTemplate(w, "pages/login.html", map[string]any{
			"Error": "Invalid email or password",
			"Email": email,
		})
		return
	}

	if !h.userRepo.VerifyPassword(user, password) {
		h.templates.ExecuteTemplate(w, "pages/login.html", map[string]any{
			"Error": "Invalid email or password",
			"Email": email,
		})
		return
	}

	session, err := h.userRepo.CreateSession(r.Context(), user.ID)
	if err != nil {
		h.templates.ExecuteTemplate(w, "pages/login.html", map[string]any{
			"Error": "Failed to create session",
			"Email": email,
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

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		h.userRepo.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
