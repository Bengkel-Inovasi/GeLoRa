package adaptersoutboundrepositoryuser

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

const path = "adapters/outbound/repository/user"

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
) portsoutboundrepository.User {
	return &sqliteImpl{
		dt:          dt,
		sqrQuestion: sqrQuestion,
		sqrDollar:   sqrDollar,
		log:         log,
	}
}

func (s *sqliteImpl) CreateUserSeed(ctx context.Context, id int64, name string, username string, passwordHash string, role domainmodel.UserRole, bio string) (err error) {
	const tag = path + "/CreateUserSeed"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":    err.Error(),
			"id":       id,
			"username": username,
		}
	}

	existsQuery, existsArgs, err := s.queryUserExistsById(id)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the existence query", logMeta(err))
		return domainmodel.ErrQueryBuilder
	}

	var count int
	if err = s.dt.QueryRow(ctx, existsQuery, existsArgs...).Scan(&count); err != nil {
		s.log.Error(ctx, tag, "Failed to check user existence", logMeta(err))
		return err
	}
	if count > 0 {
		return nil
	}

	query, args, err := s.queryCreateUserSeed(id, name, username, passwordHash, role, bio)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return domainmodel.ErrQueryBuilder
	}

	if _, err = s.dt.Exec(ctx, query, args...); err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		return err
	}

	return nil
}

func (s *sqliteImpl) CreateUser(ctx context.Context, name string, username string, passwordHash string) (id int64, err error) {
	const tag = path + "/CreateUser"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":    err.Error(),
			"name":     name,
			"username": username,
		}
	}

	query, args, err := s.queryCreateUser(name, username, passwordHash)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return 0, domainmodel.ErrQueryBuilder
	}

	res, err := s.dt.Exec(ctx, query, args...)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to execute the query", logMeta(err))
		if strings.Contains(err.Error(), "users.username") {
			return 0, domainmodel.ErrUserUsernameExists
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

func (s *sqliteImpl) ReadUserById(ctx context.Context, id int64) (user *domainmodel.User, err error) {
	const tag = path + "/ReadUserById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryReadUserById(id)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	user = &domainmodel.User{}

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *sqliteImpl) ReadUserByUsername(ctx context.Context, username string) (user *domainmodel.User, err error) {
	const tag = path + "/ReadUserByUsername"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":    err.Error(),
			"username": username,
		}
	}

	query, args, err := s.queryReadUserByUsername(username)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to build the query", logMeta(err))
		return nil, domainmodel.ErrQueryBuilder
	}

	user = &domainmodel.User{}

	row := s.dt.QueryRow(ctx, query, args...)
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		s.log.Error(ctx, tag, "Failed to scan the row", logMeta(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainmodel.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *sqliteImpl) ReadUsersByPagination(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (users []domainmodel.User, total int, err error) {
	const tag = path + "/ReadUsersByPagination"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":     err.Error(),
			"page":      page,
			"limit":     limit,
			"cursor_id": cursorId,
			"search":    search,
			"role":      role,
		}
	}

	query, args, countQuery, countArgs, err := s.queryReadUsersByPagination(page, limit, cursorId, search, role)
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
		var u domainmodel.User
		if err = rows.Scan(
			&u.Id,
			&u.Name,
			&u.Username,
			&u.PasswordHash,
			&u.Role,
			&u.Bio,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			s.log.Error(ctx, tag, "Failed to scan a row", logMeta(err))
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		s.log.Error(ctx, tag, "Row iteration error", logMeta(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (s *sqliteImpl) UpdateUserInfoById(ctx context.Context, id int64, name *string, bio *string) (err error) {
	const tag = path + "/UpdateUserInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
			"name":  name,
			"bio":   bio,
		}
	}

	if name == nil && bio == nil {
		return nil
	}

	query, args, err := s.queryUpdateUserInfoById(id, name, bio)
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
		return domainmodel.ErrUserNotFound
	}

	return nil
}

func (s *sqliteImpl) UpdateUserPasswordById(ctx context.Context, id int64, passwordHash string) (err error) {
	const tag = path + "/UpdateUserPasswordById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryUpdateUserPasswordById(id, passwordHash)
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
		return domainmodel.ErrUserNotFound
	}

	return nil
}

func (s *sqliteImpl) UpdateUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) (err error) {
	const tag = path + "/UpdateUserRoleById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
			"role":  role,
		}
	}

	query, args, err := s.queryUpdateUserRoleById(id, role)
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
		return domainmodel.ErrUserNotFound
	}

	return nil
}

func (s *sqliteImpl) DeleteUserById(ctx context.Context, id int64) (err error) {
	const tag = path + "/DeleteUserById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	query, args, err := s.queryDeleteUserById(id)
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
		return domainmodel.ErrUserNotFound
	}

	return nil
}
