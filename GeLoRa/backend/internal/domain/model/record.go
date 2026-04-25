package domainmodel

import "time"

type Record struct {
	Id              int64
	SessionId       *int64
	Time            time.Time
	HeartRate       *float64
	BodyTemperature *float64
	Latitude        *float64
	Longitude       *float64
}
