package adaptersoutboundrepositoryalert

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Masterminds/squirrel"
)

const path = "adapters/outbound/repository/alert"

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
) portsoutboundrepository.Alert {
	return &sqliteImpl{
		dt:          dt,
		sqrQuestion: sqrQuestion,
		sqrDollar:   sqrDollar,
		log:         log,
	}
}

func (s *sqliteImpl) CreateAlert(ctx context.Context, userId int64, nodeId *int64, message string) (id int64, err error) {
	const tag = path + "/CreateAlert"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{"error": err.Error(), "user_id": userId}
	}

	query, args, err := s.queryCreateAlert(userId, nodeId, message)
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
		return 0, err
	}

	return id, nil
}

func (s *sqliteImpl) ReadAlerts(ctx context.Context, unacknowledgedOnly bool) (alerts []domainmodel.Alert, err error) {
	const tag = path + "/ReadAlerts"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{"error": err.Error()}
	}

	query, args, err := s.queryReadAlerts(unacknowledgedOnly)
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
		var a domainmodel.Alert
		if err = rows.Scan(&a.Id, &a.UserId, &a.NodeId, &a.Message, &a.AcknowledgedAt, &a.CreatedAt); err != nil {
			s.log.Error(ctx, tag, "Failed to scan a row", logMeta(err))
			return nil, err
		}
		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		s.log.Error(ctx, tag, "Row iteration error", logMeta(err))
		return nil, err
	}

	return alerts, nil
}

func (s *sqliteImpl) UpdateAlertAcknowledge(ctx context.Context, id int64) (err error) {
	const tag = path + "/UpdateAlertAcknowledge"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{"error": err.Error(), "id": id}
	}

	query, args, err := s.queryUpdateAlertAcknowledge(id)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return domainmodel.ErrQueryBuilder
	}

	res, err := s.dt.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		s.log.Error(ctx, tag, "Failed to retrieve rows affected", logMeta(err))
		return err
	}
	if affected == 0 {
		return domainmodel.ErrAlertNotFound
	}

	return nil
}
