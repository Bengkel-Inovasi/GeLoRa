package adaptersinboundmqttdto

import "time"

type RecordPayload struct {
	Mid             string     `json:"mid" example:"sensor-1824721"`
	Time            *time.Time `json:"time" example:"2024-01-01T12:00:00Z"`
	HeartRate       *float64   `json:"heart_rate" example:"72.5"`
	BodyTemperature *float64   `json:"body_temperature" example:"36.6"`
	Latitude        *float64   `json:"latitude" example:"-6.2088"`
	Longitude       *float64   `json:"longitude" example:"106.8456"`
}
