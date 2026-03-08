package repository

import (
	"context"
	"errors"

	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CustomerRepo struct {
	db *database.DB
}

func NewCustomerRepo(db *database.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var c models.Customer
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, email, phone, address, city, state, zip, notes, deleted_at, created_at, updated_at
		FROM customers WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.City, &c.State, &c.Zip, &c.Notes, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) List(ctx context.Context) ([]models.Customer, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, email, phone, address, city, state, zip, notes, deleted_at, created_at, updated_at
		FROM customers WHERE deleted_at IS NULL ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.City, &c.State, &c.Zip, &c.Notes, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *CustomerRepo) Search(ctx context.Context, query string) ([]models.Customer, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, email, phone, address, city, state, zip, notes, deleted_at, created_at, updated_at
		FROM customers 
		WHERE deleted_at IS NULL 
		AND (name ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1)
		ORDER BY name
		LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.City, &c.State, &c.Zip, &c.Notes, &c.DeletedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *CustomerRepo) Create(ctx context.Context, c *models.Customer) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO customers (name, email, phone, address, city, state, zip, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, c.Name, c.Email, c.Phone, c.Address, c.City, c.State, c.Zip, c.Notes).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CustomerRepo) Update(ctx context.Context, c *models.Customer) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE customers SET 
			name = $2, email = $3, phone = $4, address = $5, 
			city = $6, state = $7, zip = $8, notes = $9, updated_at = NOW()
		WHERE id = $1
	`, c.ID, c.Name, c.Email, c.Phone, c.Address, c.City, c.State, c.Zip, c.Notes)
	return err
}

func (r *CustomerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE customers SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}
