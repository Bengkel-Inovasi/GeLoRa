package portsinboundhttp

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type User interface {
	SignIn(ctx context.Context, username string, password string) (accessToken string, refreshToken string, err error)
	SignUp(ctx context.Context, name string, username string, password string) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, id int64) (accessToken string, refreshToken string, err error)
	GetUserInfoById(ctx context.Context, id int64) (user *domainmodel.User, err error)
	GetUsersList(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (users []domainmodel.User, total int, err error)
	SetUserInfoById(ctx context.Context, id int64, name *string, bio *string) (err error)
	SetUserPasswordById(ctx context.Context, id int64, oldPassword string, newPassword string) (err error)
	SetUserPasswordDirectById(ctx context.Context, id int64, newPassword string) (err error)
	SetUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) (err error)
	RemoveUserById(ctx context.Context, id int64) (err error)
}
