package adaptersinboundhttpdto

import (
	"errors"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/utils"
)

// -- POST /auth/sign-in ------------------------------------------------------

type PostSignInRequest struct {
	Username string `json:"username" example:"john_doe_99"`
	Password string `json:"password" example:"secret123"`
}

type PostSignInResponse struct {
	AccessToken  string `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func (r *PostSignInRequest) Validate() error {
	if len(r.Username) == 0 {
		return errors.New("Username is required")
	}
	if len(r.Password) == 0 {
		return errors.New("Password is required")
	}
	return nil
}

// -- POST /auth/sign-up ------------------------------------------------------

type PostSignUpRequest struct {
	Name     string `json:"name"     example:"John Doe"`
	Username string `json:"username" example:"john_doe_99"`
	Password string `json:"password" example:"secret123"`
}

type PostSignUpResponse struct {
	AccessToken  string `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func (r *PostSignUpRequest) Validate() error {
	if err := utils.ValidateName(r.Name); err != nil {
		return err
	}
	if err := utils.ValidateUsername(r.Username); err != nil {
		return err
	}
	if err := utils.ValidatePassword(r.Password); err != nil {
		return err
	}
	return nil
}

// -- POST /auth/refresh ------------------------------------------------------

type PostRefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type PostRefreshResponse struct {
	AccessToken  string `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// -- GET /users/me -----------------------------------------------------------

type GetUserInfoMeResponse struct {
	Id        int64                `json:"id"         example:"1"`
	Name      string               `json:"name"       example:"John Doe"`
	Username  string               `json:"username"   example:"john_doe_99"`
	Role      domainmodel.UserRole `json:"role"       example:"client"`
	Bio       string               `json:"bio"        example:"Hello, I am John."`
	CreatedAt time.Time            `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt *time.Time           `json:"updated_at" example:"2024-06-01T00:00:00Z"`
}

func ConvertUserModelToGetUserInfoMeResponse(user *domainmodel.User) *GetUserInfoMeResponse {
	return &GetUserInfoMeResponse{
		Id:        user.Id,
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// -- GET /users/:id ----------------------------------------------------------

type GetUserInfoByIdResponse struct {
	Id        int64                `json:"id"         example:"1"`
	Name      string               `json:"name"       example:"John Doe"`
	Username  string               `json:"username"   example:"john_doe_99"`
	Role      domainmodel.UserRole `json:"role"       example:"client"`
	Bio       string               `json:"bio"        example:"Hello, I am John."`
	CreatedAt time.Time            `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt *time.Time           `json:"updated_at" example:"2024-06-01T00:00:00Z"`
}

func ConvertUserModelToGetUserInfoByIdResponse(user *domainmodel.User) *GetUserInfoByIdResponse {
	return &GetUserInfoByIdResponse{
		Id:        user.Id,
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// -- GET /users --------------------------------------------------------------

type GetUsersListData struct {
	Id       int64                `json:"id"       example:"1"`
	Name     string               `json:"name"     example:"John Doe"`
	Username string               `json:"username" example:"john_doe_99"`
	Role     domainmodel.UserRole `json:"role"     example:"client"`
}

type GetUsersListResponse struct {
	Data       []GetUsersListData   `json:"data"`
	Pagination CommonPaginationMeta `json:"pagination"`
}

func ConvertUserModelsToGetUsersListResponse(users []domainmodel.User, page, limit, total int) *GetUsersListResponse {
	data := make([]GetUsersListData, len(users))
	for i, u := range users {
		data[i] = GetUsersListData{
			Id:       u.Id,
			Name:     u.Name,
			Username: u.Username,
			Role:     u.Role,
		}
	}
	return &GetUsersListResponse{
		Data: data,
		Pagination: CommonPaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

// -- PATCH /users/me ---------------------------------------------------------

type PatchUserInfoMeRequest struct {
	Name *string `json:"name" example:"John Doe"`
	Bio  *string `json:"bio"  example:"Hello, I am John."`
}

func (r *PatchUserInfoMeRequest) Validate() error {
	if r.Name != nil {
		if err := utils.ValidateName(*r.Name); err != nil {
			return err
		}
	}
	if r.Bio != nil {
		if err := utils.ValidateBio(*r.Bio); err != nil {
			return err
		}
	}
	return nil
}

// -- PATCH /users/:id --------------------------------------------------------

type PatchUserInfoByIdRequest struct {
	Name *string `json:"name" example:"John Doe"`
	Bio  *string `json:"bio"  example:"Hello, I am John."`
}

func (r *PatchUserInfoByIdRequest) Validate() error {
	if r.Name != nil {
		if err := utils.ValidateName(*r.Name); err != nil {
			return err
		}
	}
	if r.Bio != nil {
		if err := utils.ValidateBio(*r.Bio); err != nil {
			return err
		}
	}
	return nil
}

// -- PUT /users/me/password --------------------------------------------------

type PutUserPasswordMeRequest struct {
	OldPassword string `json:"old_password" example:"oldsecret123"`
	NewPassword string `json:"new_password" example:"newsecret456"`
}

func (r *PutUserPasswordMeRequest) Validate() error {
	if err := utils.ValidatePassword(r.OldPassword); err != nil {
		return err
	}
	if err := utils.ValidatePassword(r.NewPassword); err != nil {
		return err
	}
	return nil
}

// -- PUT /users/:id/password -------------------------------------------------

type PutUserPasswordByIdRequest struct {
	NewPassword string `json:"new_password" example:"newsecret456"`
}

func (r *PutUserPasswordByIdRequest) Validate() error {
	if err := utils.ValidatePassword(r.NewPassword); err != nil {
		return err
	}
	return nil
}

// -- PUT /users/:id/role -----------------------------------------------------

type PutUserRoleByIdRequest struct {
	Role domainmodel.UserRole `json:"role" example:"admin"`
}

func (r *PutUserRoleByIdRequest) Validate() error {
	if err := utils.ValidateRole(string(r.Role)); err != nil {
		return err
	}
	return nil
}
