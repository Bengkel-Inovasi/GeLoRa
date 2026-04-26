package domainappcore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

const path = "domain/app/core"

type Core struct {
	infrastructure *Infrastructure
	wiring         *Wiring
}

func Run() int {
	const tag = path + "/Run"
	core := &Core{}
	ctx := context.Background()

	// Load .env
	err := config.LoadEnv()
	if err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to load .env", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	// Create infrastructures
	err = core.NewInfrastructure(ctx)
	if err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create infrastructure", nil)
		return 1
	}

	// Create wiring
	err = core.NewWiring(ctx)
	if err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create wiring", nil)
		return 1
	}

	// Delivery
	err = core.NewDelivery(ctx)
	if err != nil {
		core.wiring.log.Error(ctx, tag, "Failed to deliver", nil)
		return 1
	}

	// Start GIN Engine listening
	addr := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)
	core.wiring.log.Info(ctx, tag, "HTTP Server starting...", domainmodel.LogMeta{"addr": addr})

	err = core.infrastructure.ginEngine.Run(addr)
	if err != nil && err != http.ErrServerClosed {
		core.wiring.log.Error(ctx, tag, "Failed to start HTTP server", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	return 0
}
