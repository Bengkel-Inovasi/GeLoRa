package domainmodel

import "time"

type Session struct {
	Id        int64
	UserId    *int64
	NodeId    *int64
	StartedAt time.Time
	EndedAt   *time.Time
}
