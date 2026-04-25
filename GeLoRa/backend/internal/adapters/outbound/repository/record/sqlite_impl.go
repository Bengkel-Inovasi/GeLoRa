package adaptersoutboundrepositoryrecord

import (
	"context"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Masterminds/squirrel"
)

const path = "adapters/outbound/repository/record"

type sqliteImpl struct {
	dt          sqldt.Sqldt
	sqrQuestion *squirrel.StatementBuilderType
	sqrDollar   *squirrel.StatementBuilderType
	log         portsoutboundlogging.Generic
}

func NewSqliteImpl(
	dt sqldt.Sqldt,
	sqrQuestion *squirrel.StatementBuilderType,
	sqrDollar *squirrel.StatementBuilderType,
	log portsoutboundlogging.Generic,
) portsoutboundrepository.Record {
	return &sqliteImpl{
		dt:          dt,
		sqrQuestion: sqrQuestion,
		sqrDollar:   sqrDollar,
		log:         log,
	}
}

func (s *sqliteImpl) CreateRecord(ctx context.Context, sessionId int64, t time.Time, heartRate *float64, bodyTemperature *float64, latitude *float64, longitude *float64) (id int64, err error) {
	const tag = path + "/CreateRecord"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":      err.Error(),
			"session_id": sessionId,
			"time":       t,
		}
	}

	query, args, err := s.queryCreateRecord(sessionId, t, heartRate, bodyTemperature, latitude, longitude)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return 0, domainmodel.ErrQueryBuilder
	}

	res, err := s.dt.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		s.log.Error(ctx, tag, "Failed to retrieve last insert ID", logMeta(err))
		return id, err
	}

	return id, err
}

func (s *sqliteImpl) ReadRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (records []domainmodel.Record, err error) {
	const tag = path + "/ReadRecords"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":      err.Error(),
			"session_id": sessionId,
			"user_id":    userId,
			"node_id":    nodeId,
			"start_time": startTime,
			"end_time":   endTime,
		}
	}

	query, args, err := s.queryReadRecords(sessionId, userId, nodeId, startTime, endTime)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	rows, err := s.dt.Query(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r domainmodel.Record
		if err = rows.Scan(
			&r.Id,
			&r.SessionId,
			&r.Time,
			&r.HeartRate,
			&r.BodyTemperature,
			&r.Latitude,
			&r.Longitude,
		); err != nil {
			s.log.Error(ctx, tag, "Failed to scan a row", logMeta(err))
			return nil, err
		}
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		s.log.Error(ctx, tag, "Row iteration error", logMeta(err))
		return nil, err
	}

	return records, nil
}
