package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Address   string     `json:"address,omitempty"`
	City      string     `json:"city,omitempty"`
	State     string     `json:"state,omitempty"`
	Zip       string     `json:"zip,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Customer) FullAddress() string {
	if c.Address == "" {
		return ""
	}
	addr := c.Address
	if c.City != "" {
		addr += ", " + c.City
	}
	if c.State != "" {
		addr += ", " + c.State
	}
	if c.Zip != "" {
		addr += " " + c.Zip
	}
	return addr
}
