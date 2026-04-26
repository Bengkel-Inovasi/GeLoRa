package adaptersoutboundlogginggeneric

import (
	"context"
	"log/slog"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
)

type slogImpl struct {
	logger *slog.Logger
}

func NewSlogImpl(logger *slog.Logger) portsoutboundlogging.Generic {
	return &slogImpl{
		logger: logger,
	}
}

func (s *slogImpl) Error(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	s.logger.ErrorContext(
		ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (s *slogImpl) Warn(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	s.logger.WarnContext(
		ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (s *slogImpl) Info(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	s.logger.InfoContext(
		ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (s *slogImpl) Debug(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	s.logger.DebugContext(
		ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}
