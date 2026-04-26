package adaptersinboundmqttsubscriptions

import (
	adaptersinboundmqtthandler "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/mqtt/handler"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Subscribe(
	client mqtt.Client,
	hdlRecord *adaptersinboundmqtthandler.Record,
) (err error) {
	if tkn := client.Subscribe(config.RecordSubTopic, 0, hdlRecord.Record); tkn.Wait() && tkn.Error() != nil {
		return err
	}
	return nil
}
