package domainappapitest

import (
	"log/slog"
	"os"

	adaptersoutboundlogginggeneric "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/logging/generic"
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
