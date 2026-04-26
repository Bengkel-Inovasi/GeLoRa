package domainappcore

import (
	"fmt"
	"log/slog"
	"os"

	adaptersoutboundlogginggeneric "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/logging/generic"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
)

func logTmp(logger *slog.Logger) portsoutboundlogging.Generic {
	if logger == nil {
		return adaptersoutboundlogginggeneric.NewSlogImpl(slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		))
	}
	return adaptersoutboundlogginggeneric.NewSlogImpl(logger)
}

func getSlogLevelFromDomainLogLevel(domainLogLevel domainmodel.LogLevel) slog.Level {
	var slogLevel slog.Level
	switch domainLogLevel {
	case domainmodel.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case domainmodel.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case domainmodel.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case domainmodel.LogLevelError:
		slogLevel = slog.LevelError
	}
	return slogLevel
}

func getSQLiteConnectionPath(filePath string) string {
	return fmt.Sprintf("file:%s?mode=rw&cache=shared&_journal_mode=WAL", filePath)
}
