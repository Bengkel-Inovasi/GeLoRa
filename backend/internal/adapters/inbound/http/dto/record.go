package adaptersinboundhttpdto

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

// -- GET /records ------------------------------------------------------------

type GetRecordsData struct {
	Id              int64     `json:"id"               example:"1"`
	SessionId       *int64    `json:"session_id"       example:"42"`
	Time            time.Time `json:"time"             example:"2024-01-01T12:00:00Z"`
	HeartRate       *float64  `json:"heart_rate"       example:"72.5"`
	BodyTemperature *float64  `json:"body_temperature" example:"36.6"`
	Latitude        *float64  `json:"latitude"         example:"-6.2088"`
	Longitude       *float64  `json:"longitude"        example:"106.8456"`
}

type GetRecordsResponse struct {
	Data []GetRecordsData `json:"data"`
}

func ConvertRecordModelsToGetRecordsResponse(records []domainmodel.Record) *GetRecordsResponse {
	data := make([]GetRecordsData, len(records))
	for i, r := range records {
		data[i] = GetRecordsData{
			Id:              r.Id,
			SessionId:       r.SessionId,
			Time:            r.Time,
			HeartRate:       r.HeartRate,
			BodyTemperature: r.BodyTemperature,
			Latitude:        r.Latitude,
			Longitude:       r.Longitude,
		}
	}
	return &GetRecordsResponse{Data: data}
}
