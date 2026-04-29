package portsoutboundrepository

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Alert interface {
	CreateAlert(ctx context.Context, userId int64, nodeId *int64, message string) (id int64, err error)
	ReadAlerts(ctx context.Context, unacknowledgedOnly bool) (alerts []domainmodel.Alert, err error)
	UpdateAlertAcknowledge(ctx context.Context, id int64) (err error)
}
