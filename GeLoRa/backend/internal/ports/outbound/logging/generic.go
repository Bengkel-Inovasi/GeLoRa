package portsoutboundlogging

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Generic interface {
	Error(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Warn(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Info(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Debug(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
}
