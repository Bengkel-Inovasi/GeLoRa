package portsinboundmqtt

import (
	"context"
	"time"
)

type Record interface {
	AddRecord(ctx context.Context, mid string, time time.Time, heartRate *float64, bodyTemperature *float64, latitude *float64, longitude *float64) (id int64, err error)
}
