package domainappcore

import (
	"context"

	adaptersinboundhttphandler "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/handler"
	adaptersinboundhttproutes "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/routes"
	adaptersinboundmqtthandler "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/mqtt/handler"
	adaptersinboundmqttsubscriptions "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/mqtt/subscriptions"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

func (c *Core) NewDelivery(ctx context.Context) (err error) {
	const tag = path + "/NewDelivery"

	adaptersinboundhttproutes.Route(
		c.infrastructure.ginEngine.Group(""),
		c.wiring.log,
		c.wiring.svcHttpUser,
		adaptersinboundhttphandler.NewUser(c.wiring.svcHttpUser),
		adaptersinboundhttphandler.NewNode(c.wiring.svcHttpNode),
		adaptersinboundhttphandler.NewSession(c.wiring.svcHttpSession),
		adaptersinboundhttphandler.NewRecord(c.wiring.svcHttpRecord),
		adaptersinboundhttphandler.NewAlert(c.wiring.svcHttpAlert),
	)
	c.wiring.log.Info(ctx, tag, "HTTP Routes delivered", nil)

	err = adaptersinboundmqttsubscriptions.Subscribe(
		c.infrastructure.mqttClient,
		adaptersinboundmqtthandler.NewRecord(c.wiring.log, c.wiring.svcMqttRecord),
	)
	if err != nil {
		c.wiring.log.Error(ctx, tag, "Failed to subscribe MQTT Topics", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	c.wiring.log.Info(ctx, tag, "MQTT Subscriptions delivered", nil)

	return nil
}
