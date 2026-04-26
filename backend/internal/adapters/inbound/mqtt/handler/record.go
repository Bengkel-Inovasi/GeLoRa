package adaptersinboundmqtthandler

import (
	"context"
	"encoding/json"
	"time"

	adaptersinboundmqttdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/mqtt/dto"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundmqtt "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/mqtt"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const path = "adapters/inbound/mqtt/handler"

type Record struct {
	log       portsoutboundlogging.Generic
	svcRecord portsinboundmqtt.Record
}

func NewRecord(log portsoutboundlogging.Generic, svcRecord portsinboundmqtt.Record) *Record {
	return &Record{
		log:       log,
		svcRecord: svcRecord,
	}
}

func (r *Record) Record(client mqtt.Client, msg mqtt.Message) {
	const tag = path + "/Record"
	ctx := context.Background()

	r.log.Info(ctx, tag, "MQTT message received", domainmodel.LogMeta{
		"topic":      msg.Topic(),
		"qos":        msg.Qos(),
		"message_id": msg.MessageID(),
	})

	var payload adaptersinboundmqttdto.RecordPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		r.log.Error(ctx, tag, "Failed to parse MQTT payload", domainmodel.LogMeta{"error": err.Error()})
		return
	}

	recordTime := time.Now()
	if payload.Time != nil {
		recordTime = *payload.Time
	}

	if _, err := r.svcRecord.AddRecord(ctx, payload.Mid, recordTime, payload.HeartRate, payload.BodyTemperature, payload.Latitude, payload.Longitude); err != nil {
		r.log.Error(ctx, tag, "Failed to add record", domainmodel.LogMeta{"error": err.Error(), "mid": payload.Mid})
		return
	}
}
