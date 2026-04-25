package domainservicehttpuser

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/pwd"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/token"
)

const path = "domain/service/http/user"

type impl struct {
	tx       *sqldt.Transactor
	log      portsoutboundlogging.Generic
	repoUser portsoutboundrepository.User
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoUser portsoutboundrepository.User,
) portsinboundhttp.User {
	return &impl{
		tx:       tx,
		log:      log,
		repoUser: repoUser,
	}
}

func (i *impl) SignIn(ctx context.Context, username string, password string) (accessToken string, refreshToken string, err error) {
	const tag = path + "/SignIn"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":    err.Error(),
			"username": username,
		}
	}

	user, err := i.repoUser.ReadUserByUsername(ctx, username)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read user by username", logMeta(err))
		return "", "", err
	}

	if err = pwd.Compare(user.PasswordHash, password); err != nil {
		i.log.Error(ctx, tag, "Password mismatch", logMeta(domainmodel.ErrUserInvalidCredentials))
		return "", "", domainmodel.ErrUserInvalidCredentials
	}

	accessToken, err = token.GenerateAccess(&domainmodel.UserAccessTokenClaims{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate access token", logMeta(err))
		return "", "", err
	}

	refreshToken, err = token.GenerateRefresh(&domainmodel.UserRefreshTokenClaims{
		Id: user.Id,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate refresh token", logMeta(err))
		return "", "", err
	}

	i.log.Info(ctx, tag, "User signed in", domainmodel.LogMeta{"username": username})
	return accessToken, refreshToken, nil
}

func (i *impl) SignUp(ctx context.Context, name string, username string, password string) (accessToken string, refreshToken string, err error) {
	const tag = path + "/SignUp"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":    err.Error(),
			"name":     name,
			"username": username,
		}
	}

	passwordHash, err := pwd.Hash(password)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to hash password", logMeta(err))
		return "", "", err
	}

	id, err := i.repoUser.CreateUser(ctx, name, username, passwordHash)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to create user", logMeta(err))
		return "", "", err
	}

	accessToken, err = token.GenerateAccess(&domainmodel.UserAccessTokenClaims{
		Id:       id,
		Name:     name,
		Username: username,
		Role:     domainmodel.UserRoleClient,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate access token", logMeta(err))
		return "", "", err
	}

	refreshToken, err = token.GenerateRefresh(&domainmodel.UserRefreshTokenClaims{
		Id: id,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate refresh token", logMeta(err))
		return "", "", err
	}

	i.log.Info(ctx, tag, "User signed up", domainmodel.LogMeta{"id": id, "username": username})
	return accessToken, refreshToken, nil
}

func (i *impl) Refresh(ctx context.Context, id int64) (accessToken string, refreshToken string, err error) {
	const tag = path + "/Refresh"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	user, err := i.repoUser.ReadUserById(ctx, id)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read user by ID", logMeta(err))
		return "", "", err
	}

	accessToken, err = token.GenerateAccess(&domainmodel.UserAccessTokenClaims{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate access token", logMeta(err))
		return "", "", err
	}

	refreshToken, err = token.GenerateRefresh(&domainmodel.UserRefreshTokenClaims{
		Id: user.Id,
	})
	if err != nil {
		i.log.Error(ctx, tag, "Failed to generate refresh token", logMeta(err))
		return "", "", err
	}

	i.log.Info(ctx, tag, "Token refreshed", domainmodel.LogMeta{"id": id})
	return accessToken, refreshToken, nil
}

func (i *impl) GetUserInfoById(ctx context.Context, id int64) (user *domainmodel.User, err error) {
	const tag = path + "/GetUserInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	user, err = i.repoUser.ReadUserById(ctx, id)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read user by ID", logMeta(err))
		return nil, err
	}

	return user, nil
}

func (i *impl) GetUsersList(ctx context.Context, page int, limit int, cursorId *int64, search *string, role *domainmodel.UserRole) (users []domainmodel.User, total int, err error) {
	const tag = path + "/GetUsersList"
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

	users, total, err = i.repoUser.ReadUsersByPagination(ctx, page, limit, cursorId, search, role)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read users by pagination", logMeta(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (i *impl) SetUserInfoById(ctx context.Context, id int64, name *string, bio *string) (err error) {
	const tag = path + "/SetUserInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
			"name":  name,
			"bio":   bio,
		}
	}

	if err = i.repoUser.UpdateUserInfoById(ctx, id, name, bio); err != nil {
		i.log.Error(ctx, tag, "Failed to update user info", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "User info updated", domainmodel.LogMeta{"id": id})
	return nil
}

func (i *impl) SetUserPasswordById(ctx context.Context, id int64, oldPassword string, newPassword string) (err error) {
	const tag = path + "/SetUserPasswordById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	user, err := i.repoUser.ReadUserById(ctx, id)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read user by ID", logMeta(err))
		return err
	}

	if err = pwd.Compare(user.PasswordHash, oldPassword); err != nil {
		i.log.Error(ctx, tag, "Old password mismatch", logMeta(domainmodel.ErrUserInvalidCredentials))
		return domainmodel.ErrUserInvalidCredentials
	}

	newHash, err := pwd.Hash(newPassword)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to hash new password", logMeta(err))
		return err
	}

	if err = i.repoUser.UpdateUserPasswordById(ctx, id, newHash); err != nil {
		i.log.Error(ctx, tag, "Failed to update user password", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "User password updated", domainmodel.LogMeta{"id": id})
	return nil
}

func (i *impl) SetUserPasswordDirectById(ctx context.Context, id int64, newPassword string) (err error) {
	const tag = path + "/SetUserPasswordDirectById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	newHash, err := pwd.Hash(newPassword)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to hash new password", logMeta(err))
		return err
	}

	if err = i.repoUser.UpdateUserPasswordById(ctx, id, newHash); err != nil {
		i.log.Error(ctx, tag, "Failed to update user password", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "User password updated directly", domainmodel.LogMeta{"id": id})
	return nil
}

func (i *impl) SetUserRoleById(ctx context.Context, id int64, role domainmodel.UserRole) (err error) {
	const tag = path + "/SetUserRoleById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
			"role":  role,
		}
	}

	if err = i.repoUser.UpdateUserRoleById(ctx, id, role); err != nil {
		i.log.Error(ctx, tag, "Failed to update user role", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "User role updated", domainmodel.LogMeta{"id": id, "role": role})
	return nil
}

func (i *impl) RemoveUserById(ctx context.Context, id int64) (err error) {
	const tag = path + "/RemoveUserById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	if err = i.repoUser.DeleteUserById(ctx, id); err != nil {
		i.log.Error(ctx, tag, "Failed to delete user by ID", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "User removed", domainmodel.LogMeta{"id": id})
	return nil
}
