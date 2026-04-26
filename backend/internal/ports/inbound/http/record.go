package portsinboundhttp

import (
	"context"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Record interface {
	GetRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, end *time.Time) (records []domainmodel.Record, err error)
}
