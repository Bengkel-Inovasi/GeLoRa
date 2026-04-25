package domainappapitest

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
)

type Infrastructure struct {
	log     *slog.Logger
	baseURL string
}

func (a *ApiTest) NewInfrastructure(ctx context.Context) error {
	const tag = path + "/NewInfrastructure"

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	baseURL := fmt.Sprintf("http://localhost:%d", config.AppPort)

	a.infrastructure = &Infrastructure{
		log:     log,
		baseURL: baseURL,
	}
	logTmp(log).Info(ctx, tag, "Infrastructure created successfully", nil)
	return nil
}
