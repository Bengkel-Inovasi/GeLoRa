package domainmodel

import "time"

type Node struct {
	Id          int64
	Mid         string
	Name        string
	Description string
	IsValidated bool
	ValidatedBy *int64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
