package adaptersoutboundrepositorysession

import (
	"context"
	"database/sql"
	"errors"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Masterminds/squirrel"
)

const path = "adapters/outbound/repository/session"

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
) portsoutboundrepository.Session {
	return &sqliteImpl{
		dt:          dt,
		sqrQuestion: sqrQuestion,
		sqrDollar:   sqrDollar,
		log:         log,
	}
}

func (s *sqliteImpl) CreateSession(ctx context.Context, userId int64, nodeId int64) (id int64, err error) {
	const tag = path + "/CreateSession"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":   err.Error(),
			"user_id": userId,
			"node_id": nodeId,
		}
	}

	query, args, err := s.queryCreateSession(userId, nodeId)
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

func (s *sqliteImpl) ReadSessionById(ctx context.Context, id int64) (session *domainmodel.Session, err error) {
	const tag = path + "/ReadSessionById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryReadSessionById(id)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	session = &domainmodel.Session{}

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&session.Id,
		&session.UserId,
		&session.NodeId,
		&session.StartedAt,
		&session.EndedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrSessionNotFound
		}
		return nil, err
	}

	return session, nil
}

func (s *sqliteImpl) ReadSessionLatest(ctx context.Context, userId *int64, nodeId *int64) (session *domainmodel.Session, err error) {
	const tag = path + "/ReadSessionLatest"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":   err.Error(),
			"user_id": userId,
			"node_id": nodeId,
		}
	}

	if userId == nil && nodeId == nil {
		return nil, domainmodel.ErrSessionNoIdentifier
	}

	query, args, err := s.queryReadSessionLatest(userId, nodeId)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	session = &domainmodel.Session{}

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&session.Id,
		&session.UserId,
		&session.NodeId,
		&session.StartedAt,
		&session.EndedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrSessionNotFound
		}
		return nil, err
	}

	return session, nil
}

func (s *sqliteImpl) ReadSessionsByPagination(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (sessions []domainmodel.Session, total int, err error) {
	const tag = path + "/ReadSessionsByPagination"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":     err.Error(),
			"page":      page,
			"limit":     limit,
			"cursor_id": cursorId,
			"user_id":   userId,
			"node_id":   nodeId,
			"active":    active,
		}
	}

	query, args, countQuery, countArgs, err := s.queryReadSessionsByPagination(page, limit, cursorId, userId, nodeId, active)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, 0, domainmodel.ErrQueryBuilder
	}

	if err = s.dt.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		s.log.Error(ctx, tag, "Failed to execute the count query", logMeta(err))
		return nil, 0, err
	}

	rows, err := s.dt.Query(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var sess domainmodel.Session
		if err = rows.Scan(
			&sess.Id,
			&sess.UserId,
			&sess.NodeId,
			&sess.StartedAt,
			&sess.EndedAt,
		); err != nil {
			s.log.Error(ctx, tag, "Failed to scan a row", logMeta(err))
			return nil, 0, err
		}
		sessions = append(sessions, sess)
	}

	if err = rows.Err(); err != nil {
		s.log.Error(ctx, tag, "Row iteration error", logMeta(err))
		return nil, 0, err
	}

	return sessions, total, nil
}

func (s *sqliteImpl) UpdateSessionEndSession(ctx context.Context, id *int64, userId *int64, nodeId *int64) (err error) {
	const tag = path + "/UpdateSessionEndSession"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":   err.Error(),
			"id":      id,
			"user_id": userId,
			"node_id": nodeId,
		}
	}

	if id == nil && userId == nil && nodeId == nil {
		return domainmodel.ErrSessionNoIdentifier
	}

	query, args, err := s.queryUpdateSessionEndSession(id, userId, nodeId)
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
		return domainmodel.ErrSessionNotFound
	}

	return nil
}

func (s *sqliteImpl) DeleteSessionById(ctx context.Context, id int64) (err error) {
	const tag = path + "/DeleteSessionById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryDeleteSessionById(id)
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
		return domainmodel.ErrSessionNotFound
	}

	return nil
}
