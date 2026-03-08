package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type JobRepo struct {
	db *database.DB
}

func NewJobRepo(db *database.DB) *JobRepo {
	return &JobRepo{db: db}
}

func (r *JobRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var j models.Job
	var scheduledDate *time.Time
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, customer_id, assigned_to, title, description, status, priority,
			scheduled_date, scheduled_time, estimated_duration, completed_at,
			use_customer_address, location_address, location_city, location_state, location_zip,
			deleted_at, created_at, updated_at
		FROM jobs WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&j.ID, &j.CustomerID, &j.AssignedTo, &j.Title, &j.Description, &j.Status, &j.Priority,
		&scheduledDate, &j.ScheduledTime, &j.EstimatedDuration, &j.CompletedAt,
		&j.UseCustomerAddress, &j.LocationAddress, &j.LocationCity, &j.LocationState, &j.LocationZip,
		&j.DeletedAt, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	j.ScheduledDate = scheduledDate
	return &j, nil
}

func (r *JobRepo) List(ctx context.Context) ([]models.Job, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT j.id, j.customer_id, j.assigned_to, j.title, j.description, j.status, j.priority,
			j.scheduled_date, j.scheduled_time, j.estimated_duration, j.completed_at,
			j.use_customer_address, j.location_address, j.location_city, j.location_state, j.location_zip,
			j.deleted_at, j.created_at, j.updated_at,
			c.name as customer_name, u.name as assigned_name
		FROM jobs j
		LEFT JOIN customers c ON j.customer_id = c.id
		LEFT JOIN users u ON j.assigned_to = u.id
		WHERE j.deleted_at IS NULL 
		ORDER BY j.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var scheduledDate *time.Time
		var customerName, assignedName *string
		if err := rows.Scan(
			&j.ID, &j.CustomerID, &j.AssignedTo, &j.Title, &j.Description, &j.Status, &j.Priority,
			&scheduledDate, &j.ScheduledTime, &j.EstimatedDuration, &j.CompletedAt,
			&j.UseCustomerAddress, &j.LocationAddress, &j.LocationCity, &j.LocationState, &j.LocationZip,
			&j.DeletedAt, &j.CreatedAt, &j.UpdatedAt,
			&customerName, &assignedName,
		); err != nil {
			return nil, err
		}
		j.ScheduledDate = scheduledDate
		if customerName != nil {
			j.Customer = &models.Customer{Name: *customerName}
		}
		if assignedName != nil {
			j.AssignedUser = &models.User{Name: *assignedName}
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *JobRepo) ListByDate(ctx context.Context, date time.Time) ([]models.Job, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT j.id, j.customer_id, j.assigned_to, j.title, j.description, j.status, j.priority,
			j.scheduled_date, j.scheduled_time, j.estimated_duration, j.completed_at,
			j.use_customer_address, j.location_address, j.location_city, j.location_state, j.location_zip,
			j.deleted_at, j.created_at, j.updated_at,
			c.name as customer_name, u.name as assigned_name
		FROM jobs j
		LEFT JOIN customers c ON j.customer_id = c.id
		LEFT JOIN users u ON j.assigned_to = u.id
		WHERE j.deleted_at IS NULL AND j.scheduled_date = $1
		ORDER BY j.scheduled_time
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var scheduledDate *time.Time
		var customerName, assignedName *string
		if err := rows.Scan(
			&j.ID, &j.CustomerID, &j.AssignedTo, &j.Title, &j.Description, &j.Status, &j.Priority,
			&scheduledDate, &j.ScheduledTime, &j.EstimatedDuration, &j.CompletedAt,
			&j.UseCustomerAddress, &j.LocationAddress, &j.LocationCity, &j.LocationState, &j.LocationZip,
			&j.DeletedAt, &j.CreatedAt, &j.UpdatedAt,
			&customerName, &assignedName,
		); err != nil {
			return nil, err
		}
		j.ScheduledDate = scheduledDate
		if customerName != nil {
			j.Customer = &models.Customer{Name: *customerName}
		}
		if assignedName != nil {
			j.AssignedUser = &models.User{Name: *assignedName}
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *JobRepo) ListByWorker(ctx context.Context, workerID uuid.UUID) ([]models.Job, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT j.id, j.customer_id, j.assigned_to, j.title, j.description, j.status, j.priority,
			j.scheduled_date, j.scheduled_time, j.estimated_duration, j.completed_at,
			j.use_customer_address, j.location_address, j.location_city, j.location_state, j.location_zip,
			j.deleted_at, j.created_at, j.updated_at,
			c.name as customer_name
		FROM jobs j
		LEFT JOIN customers c ON j.customer_id = c.id
		WHERE j.deleted_at IS NULL AND j.assigned_to = $1
		ORDER BY j.scheduled_date, j.scheduled_time
	`, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var scheduledDate *time.Time
		var customerName *string
		if err := rows.Scan(
			&j.ID, &j.CustomerID, &j.AssignedTo, &j.Title, &j.Description, &j.Status, &j.Priority,
			&scheduledDate, &j.ScheduledTime, &j.EstimatedDuration, &j.CompletedAt,
			&j.UseCustomerAddress, &j.LocationAddress, &j.LocationCity, &j.LocationState, &j.LocationZip,
			&j.DeletedAt, &j.CreatedAt, &j.UpdatedAt,
			&customerName,
		); err != nil {
			return nil, err
		}
		j.ScheduledDate = scheduledDate
		if customerName != nil {
			j.Customer = &models.Customer{Name: *customerName}
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *JobRepo) Create(ctx context.Context, j *models.Job) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO jobs (customer_id, assigned_to, title, description, status, priority,
			scheduled_date, scheduled_time, estimated_duration,
			use_customer_address, location_address, location_city, location_state, location_zip)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`, j.CustomerID, j.AssignedTo, j.Title, j.Description, j.Status, j.Priority,
		j.ScheduledDate, j.ScheduledTime, j.EstimatedDuration,
		j.UseCustomerAddress, j.LocationAddress, j.LocationCity, j.LocationState, j.LocationZip,
	).Scan(&j.ID, &j.CreatedAt, &j.UpdatedAt)
}

func (r *JobRepo) Update(ctx context.Context, j *models.Job) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE jobs SET 
			customer_id = $2, assigned_to = $3, title = $4, description = $5, 
			status = $6, priority = $7, scheduled_date = $8, scheduled_time = $9, 
			estimated_duration = $10, use_customer_address = $11,
			location_address = $12, location_city = $13, location_state = $14, location_zip = $15,
			updated_at = NOW()
		WHERE id = $1
	`, j.ID, j.CustomerID, j.AssignedTo, j.Title, j.Description,
		j.Status, j.Priority, j.ScheduledDate, j.ScheduledTime,
		j.EstimatedDuration, j.UseCustomerAddress,
		j.LocationAddress, j.LocationCity, j.LocationState, j.LocationZip,
	)
	return err
}

func (r *JobRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.JobStatus, changedBy uuid.UUID) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get old status
	var oldStatus string
	err = tx.QueryRow(ctx, "SELECT status FROM jobs WHERE id = $1", id).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// Update status
	var completedAt *time.Time
	if status == models.JobStatusCompleted {
		now := time.Now()
		completedAt = &now
	}

	_, err = tx.Exec(ctx, `
		UPDATE jobs SET status = $2, completed_at = $3, updated_at = NOW() WHERE id = $1
	`, id, status, completedAt)
	if err != nil {
		return err
	}

	// Record history
	_, err = tx.Exec(ctx, `
		INSERT INTO job_history (job_id, changed_by, field, old_value, new_value)
		VALUES ($1, $2, 'status', $3, $4)
	`, id, changedBy, oldStatus, string(status))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *JobRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE jobs SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}

// Notes

func (r *JobRepo) GetNotes(ctx context.Context, jobID uuid.UUID) ([]models.JobNote, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT n.id, n.job_id, n.user_id, n.content, n.created_at, u.name
		FROM job_notes n
		JOIN users u ON n.user_id = u.id
		WHERE n.job_id = $1
		ORDER BY n.created_at DESC
	`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.JobNote
	for rows.Next() {
		var n models.JobNote
		var userName string
		if err := rows.Scan(&n.ID, &n.JobID, &n.UserID, &n.Content, &n.CreatedAt, &userName); err != nil {
			return nil, err
		}
		n.User = &models.User{Name: userName}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *JobRepo) AddNote(ctx context.Context, n *models.JobNote) error {
	return r.db.Pool.QueryRow(ctx, `
		INSERT INTO job_notes (job_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, n.JobID, n.UserID, n.Content).Scan(&n.ID, &n.CreatedAt)
}

// History

func (r *JobRepo) GetHistory(ctx context.Context, jobID uuid.UUID) ([]models.JobHistory, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT h.id, h.job_id, h.changed_by, h.field, h.old_value, h.new_value, h.changed_at, u.name
		FROM job_history h
		LEFT JOIN users u ON h.changed_by = u.id
		WHERE h.job_id = $1
		ORDER BY h.changed_at DESC
	`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.JobHistory
	for rows.Next() {
		var h models.JobHistory
		var userName *string
		if err := rows.Scan(&h.ID, &h.JobID, &h.ChangedBy, &h.Field, &h.OldValue, &h.NewValue, &h.ChangedAt, &userName); err != nil {
			return nil, err
		}
		if userName != nil {
			h.User = &models.User{Name: *userName}
		}
		history = append(history, h)
	}
	return history, nil
}
