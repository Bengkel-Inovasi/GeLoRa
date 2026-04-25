package adaptersoutboundrepositoryrecord

import (
	"time"

	"github.com/Masterminds/squirrel"
)

func (s *sqliteImpl) queryCreateRecord(sessionId int64, t time.Time, heartRate *float64, bodyTemperature *float64, latitude *float64, longitude *float64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_records").
		Columns(
			"session_id",
			"time",
			"heart_rate",
			"body_temperature",
			"latitude",
			"longitude",
		).
		Values(
			sessionId,
			t,
			heartRate,
			bodyTemperature,
			latitude,
			longitude,
		).
		ToSql()
}

func (s *sqliteImpl) queryReadRecords(sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Select(
			"r.id",
			"r.session_id",
			"r.time",
			"r.heart_rate",
			"r.body_temperature",
			"r.latitude",
			"r.longitude",
		).
		From("gelora_records r").
		OrderBy("r.time ASC")

	if userId != nil || nodeId != nil {
		q = q.Join("gelora_sessions s ON r.session_id = s.id")
		if userId != nil {
			q = q.Where(squirrel.Eq{"s.user_id": *userId})
		}
		if nodeId != nil {
			q = q.Where(squirrel.Eq{"s.node_id": *nodeId})
		}
	}

	if sessionId != nil {
		q = q.Where(squirrel.Eq{"r.session_id": *sessionId})
	}

	if startTime != nil {
		q = q.Where(squirrel.GtOrEq{"r.time": *startTime})
	}

	if endTime != nil {
		q = q.Where(squirrel.LtOrEq{"r.time": *endTime})
	}

	return q.ToSql()
}
