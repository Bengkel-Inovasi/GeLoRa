package portsinboundhttp

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Alert interface {
	AddAlert(ctx context.Context, userId int64, message string) (id int64, err error)
	GetAlerts(ctx context.Context, unacknowledgedOnly bool) (alerts []domainmodel.Alert, err error)
	AcknowledgeAlert(ctx context.Context, id int64) (err error)
}
