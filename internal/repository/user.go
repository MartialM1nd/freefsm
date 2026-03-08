package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name, phone, role, deleted_at, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Phone, &u.Role, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name, phone, role, deleted_at, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Phone, &u.Role, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, email, password_hash, name, phone, role, deleted_at, created_at, updated_at
		FROM users WHERE deleted_at IS NULL ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Phone, &u.Role, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) ListTechnicians(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, email, password_hash, name, phone, role, deleted_at, created_at, updated_at
		FROM users WHERE deleted_at IS NULL AND role = 'technician' ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Phone, &u.Role, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, u *models.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, name, phone, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, u.Email, string(hash), u.Name, u.Phone, u.Role).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users SET email = $2, name = $3, phone = $4, role = $5, updated_at = NOW()
		WHERE id = $1
	`, u.ID, u.Email, u.Name, u.Phone, u.Role)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, `
		UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1
	`, id, string(hash))
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *UserRepo) VerifyPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

// Session management

func (r *UserRepo) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	s := &models.Session{
		UserID:    userID,
		Token:     hex.EncodeToString(token),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, s.UserID, s.Token, s.ExpiresAt).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *UserRepo) GetSession(ctx context.Context, token string) (*models.Session, error) {
	var s models.Session
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions WHERE token = $1 AND expires_at > NOW()
	`, token).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *UserRepo) DeleteSession(ctx context.Context, token string) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (r *UserRepo) CleanExpiredSessions(ctx context.Context) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}
