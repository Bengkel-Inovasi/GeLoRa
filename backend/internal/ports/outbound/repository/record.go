package portsoutboundrepository

import (
	"context"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Record interface {
	CreateRecord(ctx context.Context, sessionId int64, time time.Time, heartRate *float64, bodyTemperature *float64, latitude *float64, longitude *float64) (id int64, err error)
	ReadRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (records []domainmodel.Record, err error)
}
