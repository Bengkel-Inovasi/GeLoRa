package adaptersoutboundrepositorysession

import "github.com/Masterminds/squirrel"

func (s *sqliteImpl) queryCreateSession(userId int64, nodeId int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_sessions").
		Columns(
			"user_id",
			"node_id",
		).
		Values(
			userId,
			nodeId,
		).
		ToSql()
}

func (s *sqliteImpl) queryReadSessionById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select(
			"id",
			"user_id",
			"node_id",
			"started_at",
			"ended_at",
		).
		From("gelora_sessions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryReadSessionLatest(userId *int64, nodeId *int64) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Select(
			"id",
			"user_id",
			"node_id",
			"started_at",
			"ended_at",
		).
		From("gelora_sessions").
		OrderBy("started_at DESC").
		Limit(1)

	if userId != nil {
		q = q.Where(squirrel.Eq{"user_id": *userId})
	}

	if nodeId != nil {
		q = q.Where(squirrel.Eq{"node_id": *nodeId})
	}

	return q.ToSql()
}

func (s *sqliteImpl) queryReadSessionsByPagination(page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (query string, args []any, countQuery string, countArgs []any, err error) {
	q := s.sqrQuestion.
		Select(
			"id",
			"user_id",
			"node_id",
			"started_at",
			"ended_at",
		).
		From("gelora_sessions").
		OrderBy("id ASC").
		Limit(uint64(limit))

	c := s.sqrQuestion.
		Select("COUNT(*)").
		From("gelora_sessions")

	if cursorId != nil {
		q = q.Where(squirrel.Gt{"id": *cursorId})
	} else {
		q = q.Offset(uint64((page - 1) * limit))
	}

	if userId != nil {
		q = q.Where(squirrel.Eq{"user_id": *userId})
		c = c.Where(squirrel.Eq{"user_id": *userId})
	}

	if nodeId != nil {
		q = q.Where(squirrel.Eq{"node_id": *nodeId})
		c = c.Where(squirrel.Eq{"node_id": *nodeId})
	}

	if active != nil {
		if *active {
			q = q.Where(squirrel.Eq{"ended_at": nil})
			c = c.Where(squirrel.Eq{"ended_at": nil})
		} else {
			q = q.Where(squirrel.NotEq{"ended_at": nil})
			c = c.Where(squirrel.NotEq{"ended_at": nil})
		}
	}

	query, args, err = q.ToSql()
	if err != nil {
		return
	}

	countQuery, countArgs, err = c.ToSql()
	return
}

func (s *sqliteImpl) queryUpdateSessionEndSession(id *int64, userId *int64, nodeId *int64) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Update("gelora_sessions").
		Set("ended_at", squirrel.Expr("CURRENT_TIMESTAMP"))

	if id != nil {
		q = q.Where(squirrel.Eq{"id": *id})
	} else {
		if userId != nil {
			q = q.Where(squirrel.Eq{"user_id": *userId})
		}
		if nodeId != nil {
			q = q.Where(squirrel.Eq{"node_id": *nodeId})
		}
	}

	return q.ToSql()
}

func (s *sqliteImpl) queryDeleteSessionById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Delete("gelora_sessions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}
