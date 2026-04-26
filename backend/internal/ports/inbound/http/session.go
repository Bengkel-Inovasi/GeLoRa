package portsinboundhttp

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Session interface {
	AddSession(ctx context.Context, userId int64, nodeId int64) (id int64, err error)
	GetSessionById(ctx context.Context, id int64) (session *domainmodel.Session, err error)
	GetSessionsList(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (sessions []domainmodel.Session, total int, err error)
	SetSessionEndSession(ctx context.Context, id *int64, userId *int64, nodeId *int64) (err error)
	RemoveSessionById(ctx context.Context, id int64) (err error)
}
