package portsoutboundrepository

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Session interface {
	CreateSession(ctx context.Context, userId int64, nodeId int64) (id int64, err error)
	ReadSessionById(ctx context.Context, id int64) (session *domainmodel.Session, err error)
	ReadSessionLatest(ctx context.Context, userId *int64, nodeId *int64) (session *domainmodel.Session, err error)
	ReadSessionsByPagination(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (sessions []domainmodel.Session, total int, err error)
	UpdateSessionEndSession(ctx context.Context, id *int64, userId *int64, nodeId *int64) (err error)
	DeleteSessionById(ctx context.Context, id int64) (err error)
}
