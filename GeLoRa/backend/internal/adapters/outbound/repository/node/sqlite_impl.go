package adaptersoutboundrepositorynode

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Masterminds/squirrel"
)

const path = "adapters/outbound/repository/node"

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
) portsoutboundrepository.Node {
	return &sqliteImpl{
		dt:          dt,
		sqrQuestion: sqrQuestion,
		sqrDollar:   sqrDollar,
		log:         log,
	}
}

func (s *sqliteImpl) CreateNode(ctx context.Context, mid string, name string) (id int64, err error) {
	const tag = path + "/CreateNode"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"mid":   mid,
			"name":  name,
		}
	}

	query, args, err := s.queryCreateNode(mid, name)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return 0, domainmodel.ErrQueryBuilder
	}

	res, err := s.dt.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		if strings.Contains(err.Error(), "nodes.mid") {
			return 0, domainmodel.ErrNodeMidExists
		}
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		s.log.Error(ctx, tag, "Failed to retrieve last insert ID", logMeta(err))
		return id, err
	}

	return id, err
}

func (s *sqliteImpl) ReadNodeById(ctx context.Context, id int64) (node *domainmodel.Node, err error) {
	const tag = path + "/ReadNodeById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryReadNodeById(id)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	node = &domainmodel.Node{}
	var isValidated int

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&node.Id,
		&node.Mid,
		&node.Name,
		&node.Description,
		&isValidated,
		&node.ValidatedBy,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrNodeNotFound
		}
		return nil, err
	}
	node.IsValidated = isValidated != 0

	return node, nil
}

func (s *sqliteImpl) ReadNodeByMid(ctx context.Context, mid string) (node *domainmodel.Node, err error) {
	const tag = path + "/ReadNodeByMid"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"mid":   mid,
		}
	}

	query, args, err := s.queryReadNodeByMid(mid)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	node = &domainmodel.Node{}
	var isValidated int

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&node.Id,
		&node.Mid,
		&node.Name,
		&node.Description,
		&isValidated,
		&node.ValidatedBy,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrNodeNotFound
		}
		return nil, err
	}
	node.IsValidated = isValidated != 0

	return node, nil
}

func (s *sqliteImpl) ReadNodesByPagination(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (nodes []domainmodel.Node, total int, err error) {
	const tag = path + "/ReadNodesByPagination"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":        err.Error(),
			"page":         page,
			"limit":        limit,
			"cursor_id":    cursorId,
			"search":       search,
			"is_validated": isValidated,
		}
	}

	query, args, countQuery, countArgs, err := s.queryReadNodesByPagination(page, limit, cursorId, search, isValidated)
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
		var n domainmodel.Node
		var isValidatedInt int
		if err = rows.Scan(
			&n.Id,
			&n.Mid,
			&n.Name,
			&n.Description,
			&isValidatedInt,
			&n.ValidatedBy,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			s.log.Error(ctx, tag, "Failed to scan a row", logMeta(err))
			return nil, 0, err
		}
		n.IsValidated = isValidatedInt != 0
		nodes = append(nodes, n)
	}

	if err = rows.Err(); err != nil {
		s.log.Error(ctx, tag, "Row iteration error", logMeta(err))
		return nil, 0, err
	}

	return nodes, total, nil
}

func (s *sqliteImpl) UpdateNodeInfoById(ctx context.Context, id int64, name *string, description *string) (err error) {
	const tag = path + "/UpdateNodeInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":       err.Error(),
			"id":          id,
			"name":        name,
			"description": description,
		}
	}

	query, args, err := s.queryUpdateNodeInfoById(id, name, description)
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
		return domainmodel.ErrNodeNotFound
	}

	return nil
}

func (s *sqliteImpl) UpdateNodeValidation(ctx context.Context, id int64, validation bool, validatedBy int64) (err error) {
	const tag = path + "/UpdateNodeValidation"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":        err.Error(),
			"id":           id,
			"validation":   validation,
			"validated_by": validatedBy,
		}
	}

	query, args, err := s.queryUpdateNodeValidation(id, validation, validatedBy)
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
		return domainmodel.ErrNodeNotFound
	}

	return nil
}

func (s *sqliteImpl) DeleteNodeById(ctx context.Context, id int64) (err error) {
	const tag = path + "/DeleteNodeById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryDeleteNodeById(id)
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
		return domainmodel.ErrNodeNotFound
	}

	return nil
}
