package domainmodel

import "time"

type Alert struct {
	Id             int64
	UserId         *int64
	Message        string
	AcknowledgedAt *time.Time
	CreatedAt      time.Time
}
