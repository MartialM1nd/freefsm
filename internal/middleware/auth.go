package middleware

import (
	"context"
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/MartialM1nd/freefsm/internal/repository"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(db *database.DB) func(http.Handler) http.Handler {
	userRepo := repository.NewUserRepo(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			session, err := userRepo.GetSession(r.Context(), cookie.Value)
			if err != nil || session == nil {
				http.SetCookie(w, &http.Cookie{
					Name:   "session",
					Value:  "",
					MaxAge: -1,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := userRepo.GetByID(r.Context(), session.UserID)
			if err != nil || user == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
