package domainmodel

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

const (
	UserRoleSuper  UserRole = "super"
	UserRoleAdmin  UserRole = "admin"
	UserRoleClient UserRole = "client"
)

func (u UserRole) Order() int {
	switch u {
	case UserRoleSuper:
		return 3
	case UserRoleAdmin:
		return 2
	case UserRoleClient:
		return 1
	default:
		return 0
	}
}

type User struct {
	Id           int64
	Name         string
	Username     string
	PasswordHash string
	Role         UserRole
	Bio          string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type UserAccessTokenClaims struct {
	Id       int64
	Name     string
	Username string
	Role     UserRole
	jwt.RegisteredClaims
}

type UserRefreshTokenClaims struct {
	Id int64
	jwt.RegisteredClaims
}
