package adaptersoutboundrepositorynode

import "github.com/Masterminds/squirrel"

func (s *sqliteImpl) queryCreateNode(mid string, name string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_nodes").
		Columns(
			"mid",
			"name",
		).
		Values(
			mid,
			name,
		).
		ToSql()
}

func (s *sqliteImpl) queryReadNodeById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select(
			"id",
			"mid",
			"name",
			"description",
			"is_validated",
			"validated_by",
			"created_at",
			"updated_at",
		).
		From("gelora_nodes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryReadNodeByMid(mid string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select(
			"id",
			"mid",
			"name",
			"description",
			"is_validated",
			"validated_by",
			"created_at",
			"updated_at",
		).
		From("gelora_nodes").
		Where(squirrel.Eq{"mid": mid}).
		ToSql()
}

func (s *sqliteImpl) queryReadNodesByPagination(page int, limit int, cursorId *int64, search *string, isValidated *bool) (query string, args []any, countQuery string, countArgs []any, err error) {
	q := s.sqrQuestion.
		Select(
			"id",
			"mid",
			"name",
			"description",
			"is_validated",
			"validated_by",
			"created_at",
			"updated_at",
		).
		From("gelora_nodes").
		OrderBy("id ASC").
		Limit(uint64(limit))

	c := s.sqrQuestion.
		Select("COUNT(*)").
		From("gelora_nodes")

	if cursorId != nil {
		q = q.Where(squirrel.Gt{"id": *cursorId})
	} else {
		q = q.Offset(uint64((page - 1) * limit))
	}

	if search != nil {
		q = q.Where(squirrel.Or{
			squirrel.ILike{"mid": "%" + *search + "%"},
			squirrel.ILike{"name": "%" + *search + "%"},
		})
		c = c.Where(squirrel.Or{
			squirrel.ILike{"mid": "%" + *search + "%"},
			squirrel.ILike{"name": "%" + *search + "%"},
		})
	}

	if isValidated != nil {
		val := 0
		if *isValidated {
			val = 1
		}
		q = q.Where(squirrel.Eq{"is_validated": val})
		c = c.Where(squirrel.Eq{"is_validated": val})
	}

	query, args, err = q.ToSql()
	if err != nil {
		return
	}

	countQuery, countArgs, err = c.ToSql()
	return
}

func (s *sqliteImpl) queryUpdateNodeInfoById(id int64, name *string, description *string) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Update("gelora_nodes")

	if name != nil {
		q = q.Set("name", name)
	}

	if description != nil {
		q = q.Set("description", description)
	}

	return q.
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryUpdateNodeValidation(id int64, validation bool, validatedBy int64) (query string, args []any, err error) {
	isValidated := 0
	if validation {
		isValidated = 1
	}
	return s.sqrQuestion.
		Update("gelora_nodes").
		Set("is_validated", isValidated).
		Set("validated_by", validatedBy).
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryDeleteNodeById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Delete("gelora_nodes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}
