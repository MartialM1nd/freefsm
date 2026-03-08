package models

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusNew            JobStatus = "new"
	JobStatusInTransit      JobStatus = "in_transit"
	JobStatusInProgress     JobStatus = "in_progress"
	JobStatusPending        JobStatus = "pending"
	JobStatusScheduledReturn JobStatus = "scheduled_return"
	JobStatusReadyToInvoice JobStatus = "ready_to_invoice"
	JobStatusCompleted      JobStatus = "completed"
	JobStatusCancelled      JobStatus = "cancelled"
)

type JobPriority string

const (
	JobPriorityLow    JobPriority = "low"
	JobPriorityMedium JobPriority = "medium"
	JobPriorityHigh   JobPriority = "high"
	JobPriorityUrgent JobPriority = "urgent"
)

type Job struct {
	ID                 uuid.UUID   `json:"id"`
	CustomerID         *uuid.UUID  `json:"customer_id,omitempty"`
	AssignedTo         *uuid.UUID  `json:"assigned_to,omitempty"`
	Title              string      `json:"title"`
	Description        string      `json:"description,omitempty"`
	Status             JobStatus   `json:"status"`
	Priority           JobPriority `json:"priority"`
	ScheduledDate      *time.Time  `json:"scheduled_date,omitempty"`
	ScheduledTime      string      `json:"scheduled_time,omitempty"`
	EstimatedDuration  *int        `json:"estimated_duration,omitempty"` // minutes
	CompletedAt        *time.Time  `json:"completed_at,omitempty"`
	UseCustomerAddress bool        `json:"use_customer_address"`
	LocationAddress    string      `json:"location_address,omitempty"`
	LocationCity       string      `json:"location_city,omitempty"`
	LocationState      string      `json:"location_state,omitempty"`
	LocationZip        string      `json:"location_zip,omitempty"`
	DeletedAt          *time.Time  `json:"deleted_at,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`

	// Joined fields (not stored)
	Customer     *Customer `json:"customer,omitempty"`
	AssignedUser *User     `json:"assigned_user,omitempty"`
}

func (j *Job) FullLocation() string {
	if j.LocationAddress == "" {
		return ""
	}
	addr := j.LocationAddress
	if j.LocationCity != "" {
		addr += ", " + j.LocationCity
	}
	if j.LocationState != "" {
		addr += ", " + j.LocationState
	}
	if j.LocationZip != "" {
		addr += " " + j.LocationZip
	}
	return addr
}

func (j *Job) StatusLabel() string {
	labels := map[JobStatus]string{
		JobStatusNew:             "New",
		JobStatusInTransit:       "In Transit",
		JobStatusInProgress:      "In Progress",
		JobStatusPending:         "Pending",
		JobStatusScheduledReturn: "Scheduled Return",
		JobStatusReadyToInvoice:  "Ready to Invoice",
		JobStatusCompleted:       "Completed",
		JobStatusCancelled:       "Cancelled",
	}
	return labels[j.Status]
}

func (j *Job) PriorityLabel() string {
	labels := map[JobPriority]string{
		JobPriorityLow:    "Low",
		JobPriorityMedium: "Medium",
		JobPriorityHigh:   "High",
		JobPriorityUrgent: "Urgent",
	}
	return labels[j.Priority]
}

type JobNote struct {
	ID        uuid.UUID `json:"id"`
	JobID     uuid.UUID `json:"job_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Joined
	User *User `json:"user,omitempty"`
}

type JobHistory struct {
	ID        uuid.UUID  `json:"id"`
	JobID     uuid.UUID  `json:"job_id"`
	ChangedBy *uuid.UUID `json:"changed_by,omitempty"`
	Field     string     `json:"field"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
	ChangedAt time.Time  `json:"changed_at"`

	// Joined
	User *User `json:"user,omitempty"`
}
