package portsoutboundclient

import (
	"context"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Appself interface {
	SetAccessToken(accessToken string)
	SetRefreshToken(refreshToken string)

	PostAuthSignIn(ctx context.Context, username string, password string) (resp *domainmodel.AppselfAuthResponse, err error)
	PostAuthSignUp(ctx context.Context, name string, username string, password string) (resp *domainmodel.AppselfAuthResponse, err error)
	PostAuthRefresh(ctx context.Context, refreshToken string) (resp *domainmodel.AppselfAuthResponse, err error)

	GetUserInfoMe(ctx context.Context) (resp *domainmodel.AppselfUserInfo, err error)
	PatchUserInfoMe(ctx context.Context, name *string, bio *string) (err error)
	PutUserPasswordMe(ctx context.Context, oldPassword string, newPassword string) (err error)
	GetUsersList(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (resp *domainmodel.AppselfUsersList, err error)
	GetUserInfoById(ctx context.Context, id int64) (resp *domainmodel.AppselfUserInfo, err error)
	PatchUserInfoById(ctx context.Context, id int64, name *string, bio *string) (err error)
	PutUserPasswordById(ctx context.Context, id int64, newPassword string) (err error)
	PutUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) (err error)
	DeleteUserById(ctx context.Context, id int64) (err error)

	GetNodesList(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (resp *domainmodel.AppselfNodesList, err error)
	GetNodeInfoById(ctx context.Context, id int64) (resp *domainmodel.AppselfNodeInfo, err error)
	PatchNodeInfoById(ctx context.Context, id int64, name *string, description *string) (err error)
	PutNodeValidate(ctx context.Context, id int64, validatedBy int64) (err error)
	DeleteNodeById(ctx context.Context, id int64) (err error)

	PostSession(ctx context.Context, userId int64, nodeId int64) (resp *domainmodel.AppselfPostSessionResponse, err error)
	GetSessionsList(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (resp *domainmodel.AppselfSessionsList, err error)
	GetSessionById(ctx context.Context, id int64) (resp *domainmodel.AppselfSessionInfo, err error)
	PutSessionEnd(ctx context.Context, id *int64, userId *int64, nodeId *int64) (err error)
	DeleteSessionById(ctx context.Context, id int64) (err error)

	GetRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (resp *domainmodel.AppselfRecordsList, err error)
}
