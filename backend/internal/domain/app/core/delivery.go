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
		c.wiring.Log,
		c.wiring.SvcHttpUser,
		adaptersinboundhttphandler.NewUser(c.wiring.SvcHttpUser),
		adaptersinboundhttphandler.NewNode(c.wiring.SvcHttpNode),
		adaptersinboundhttphandler.NewSession(c.wiring.SvcHttpSession),
		adaptersinboundhttphandler.NewRecord(c.wiring.SvcHttpRecord),
		adaptersinboundhttphandler.NewAlert(c.wiring.SvcHttpAlert),
	)
	c.wiring.Log.Info(ctx, tag, "HTTP Routes delivered", nil)

	err = adaptersinboundmqttsubscriptions.Subscribe(
		c.infrastructure.mqttClient,
		adaptersinboundmqtthandler.NewRecord(c.wiring.Log, c.wiring.SvcMqttRecord),
	)
	if err != nil {
		c.wiring.Log.Error(ctx, tag, "Failed to subscribe MQTT Topics", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	c.wiring.Log.Info(ctx, tag, "MQTT Subscriptions delivered", nil)

	return nil
}
