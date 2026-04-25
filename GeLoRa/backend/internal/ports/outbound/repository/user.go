package portsoutboundrepository

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type User interface {
	CreateUserSeed(ctx context.Context, id int64, name string, username string, passwordHash string, role domainmodel.UserRole, bio string) (err error)
	CreateUser(ctx context.Context, name string, username string, passwordHash string) (id int64, err error)
	ReadUserById(ctx context.Context, id int64) (user *domainmodel.User, err error)
	ReadUserByUsername(ctx context.Context, username string) (user *domainmodel.User, err error)
	ReadUsersByPagination(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (users []domainmodel.User, total int, err error)
	UpdateUserInfoById(ctx context.Context, id int64, name *string, bio *string) (err error)
	UpdateUserPasswordById(ctx context.Context, id int64, passwordHash string) (err error)
	UpdateUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) (err error)
	DeleteUserById(ctx context.Context, id int64) (err error)
}
