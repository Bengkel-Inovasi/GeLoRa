package adaptersoutboundrepositoryuser

import (
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Masterminds/squirrel"
)

func (s *sqliteImpl) queryUserExistsById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select("COUNT(*)").
		From("gelora_users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryCreateUserSeed(id int64, name string, username string, passwordHash string, role domainmodel.UserRole, bio string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_users").
		Columns("id", "name", "username", "password_hash", "role", "bio").
		Values(id, name, username, passwordHash, string(role), bio).
		ToSql()
}

func (s *sqliteImpl) queryCreateUser(name string, username string, passwordHash string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Insert("gelora_users").
		Columns(
			"name",
			"username",
			"password_hash",
			"role",
		).
		Values(
			name,
			username,
			passwordHash,
			"client",
		).
		ToSql()
}

func (s *sqliteImpl) queryReadUserById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select(
			"id",
			"name",
			"username",
			"password_hash",
			"role",
			"bio",
			"created_at",
			"updated_at",
		).
		From("gelora_users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryReadUserByUsername(username string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Select(
			"id",
			"name",
			"username",
			"password_hash",
			"role",
			"bio",
			"created_at",
			"updated_at",
		).
		From("gelora_users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
}

func (s *sqliteImpl) queryReadUsersByPagination(page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (query string, args []any, countQuery string, countArgs []any, err error) {
	q := s.sqrQuestion.
		Select(
			"id",
			"name",
			"username",
			"password_hash",
			"role",
			"bio",
			"created_at",
			"updated_at",
		).
		From("gelora_users").
		OrderBy("id ASC").
		Limit(uint64(limit))

	c := s.sqrQuestion.
		Select("COUNT(*)").
		From("gelora_users")

	if cursorId != nil {
		q = q.Where(squirrel.Gt{"id": *cursorId})
	} else {
		q = q.Offset(uint64((page - 1) * limit))
	}

	if search != nil {
		q = q.Where(squirrel.ILike{"name": "%" + *search + "%"})
		c = c.Where(squirrel.ILike{"name": "%" + *search + "%"})
	}

	if role != nil {
		q = q.Where(squirrel.Eq{"role": *role})
		c = c.Where(squirrel.Eq{"role": *role})
	}

	query, args, err = q.ToSql()
	if err != nil {
		return
	}

	countQuery, countArgs, err = c.ToSql()
	return
}

func (s *sqliteImpl) queryUpdateUserInfoById(id int64, name *string, bio *string) (query string, args []any, err error) {
	q := s.sqrQuestion.
		Update("gelora_users")

	if name != nil {
		q = q.Set("name", name)
	}

	if bio != nil {
		q = q.Set("bio", bio)
	}

	return q.
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryUpdateUserPasswordById(id int64, passwordHash string) (query string, args []any, err error) {
	return s.sqrQuestion.
		Update("gelora_users").
		Set("password_hash", passwordHash).
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryUpdateUserRoleById(id int64, role domainmodel.UserRole) (query string, args []any, err error) {
	return s.sqrQuestion.
		Update("gelora_users").
		Set("role", role).
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func (s *sqliteImpl) queryDeleteUserById(id int64) (query string, args []any, err error) {
	return s.sqrQuestion.
		Delete("gelora_users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
}
