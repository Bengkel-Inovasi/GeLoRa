package domainappapitest

import (
	"context"

	adaptersoutboundclientappself "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/client/appself"
	adaptersoutboundlogginggeneric "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/logging/generic"
	portsoutboundclient "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/client"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
)

type Wiring struct {
	log     portsoutboundlogging.Generic
	appself portsoutboundclient.Appself
}

func (a *ApiTest) NewWiring(ctx context.Context) error {
	const tag = path + "/NewWiring"

	log := adaptersoutboundlogginggeneric.NewSlogImpl(a.infrastructure.log)

	appself := adaptersoutboundclientappself.New(a.infrastructure.baseURL)

	a.wiring = &Wiring{
		log:     log,
		appself: appself,
	}
	log.Info(ctx, tag, "Wiring created successfully", nil)
	return nil
}
