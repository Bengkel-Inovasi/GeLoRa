package adaptersinboundhttphandler

import (
	"errors"
	"net/http"
	"strconv"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	middlewareauth "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/auth"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	"github.com/gin-gonic/gin"
)

type User struct {
	svcUser portsinboundhttp.User
}

func NewUser(svcUser portsinboundhttp.User) *User {
	return &User{
		svcUser: svcUser,
	}
}

// PostSignIn godoc
// @Summary      Sign in
// @Description  Authenticate a user and return access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      adaptersinboundhttpdto.PostSignInRequest   true  "Credentials"
// @Success      200   {object}  adaptersinboundhttpdto.PostSignInResponse
// @Failure      400   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /auth/sign-in [post]
func (u *User) PostSignIn(c *gin.Context) {
	var req adaptersinboundhttpdto.PostSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := u.svcUser.SignIn(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound),
			errors.Is(err, domainmodel.ErrUserInvalidCredentials):
			c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid username or password"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.PostSignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// PostSignUp godoc
// @Summary      Sign up
// @Description  Register a new user and return access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      adaptersinboundhttpdto.PostSignUpRequest   true  "Registration details"
// @Success      201   {object}  adaptersinboundhttpdto.PostSignUpResponse
// @Failure      400   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      409   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /auth/sign-up [post]
func (u *User) PostSignUp(c *gin.Context) {
	var req adaptersinboundhttpdto.PostSignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := u.svcUser.SignUp(c.Request.Context(), req.Name, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserUsernameExists):
			c.JSON(http.StatusConflict, adaptersinboundhttpdto.CommonErrorResponse{Error: "Username is already taken"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusCreated, adaptersinboundhttpdto.PostSignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// PostRefresh godoc
// @Summary      Refresh tokens
// @Description  Issue a new access and refresh token pair using a valid refresh token
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  adaptersinboundhttpdto.PostRefreshResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /auth/refresh [post]
func (u *User) PostRefresh(c *gin.Context) {
	claimsAny, exists := c.Get("refresh_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}
	claims, ok := claimsAny.(*domainmodel.UserRefreshTokenClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	accessToken, refreshToken, err := u.svcUser.Refresh(c.Request.Context(), claims.Id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Account no longer exists"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.PostRefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// GetUserInfoMe godoc
// @Summary      Get my profile
// @Description  Return the profile of the currently authenticated user
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  adaptersinboundhttpdto.GetUserInfoMeResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/me [get]
func (u *User) GetUserInfoMe(c *gin.Context) {
	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	user, err := u.svcUser.GetUserInfoById(c.Request.Context(), claims.Id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertUserModelToGetUserInfoMeResponse(user))
}

// GetUserInfoById godoc
// @Summary      Get user by ID
// @Description  Return the profile of a specific user
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  adaptersinboundhttpdto.GetUserInfoByIdResponse
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/{id} [get]
func (u *User) GetUserInfoById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	user, err := u.svcUser.GetUserInfoById(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertUserModelToGetUserInfoByIdResponse(user))
}

// GetUsersList godoc
// @Summary      List users
// @Description  Return a paginated list of users with optional filters
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number"           default(1)
// @Param        limit      query     int     false  "Items per page"        default(10)
// @Param        cursor_id  query     int     false  "Cursor ID for pagination"
// @Param        search     query     string  false  "Search by name or username"
// @Param        role       query     string  false  "Filter by role (super, admin, client)"
// @Success      200        {object}  adaptersinboundhttpdto.GetUsersListResponse
// @Failure      400        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users [get]
func (u *User) GetUsersList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var cursorId *int64
	if v := c.Query("cursor_id"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid cursor_id"})
			return
		}
		cursorId = &n
	}

	var search *string
	if v := c.Query("search"); v != "" {
		search = &v
	}

	var role *domainmodel.UserRole
	if v := c.Query("role"); v != "" {
		r := domainmodel.UserRole(v)
		role = &r
	}

	users, total, err := u.svcUser.GetUsersList(c.Request.Context(), page, limit, cursorId, search, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertUserModelsToGetUsersListResponse(users, page, limit, total))
}

// PatchUserInfoMe godoc
// @Summary      Update my profile
// @Description  Update the name and/or bio of the currently authenticated user
// @Tags         Users
// @Accept       json
// @Security     BearerAuth
// @Param        body  body  adaptersinboundhttpdto.PatchUserInfoMeRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/me [patch]
func (u *User) PatchUserInfoMe(c *gin.Context) {
	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	var req adaptersinboundhttpdto.PatchUserInfoMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := u.svcUser.SetUserInfoById(c.Request.Context(), claims.Id, req.Name, req.Bio); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PatchUserInfoById godoc
// @Summary      Update user by ID
// @Description  Update the name and/or bio of a specific user
// @Tags         Users
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  int                           true  "User ID"
// @Param        body  body  adaptersinboundhttpdto.PatchUserInfoByIdRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/{id} [patch]
func (u *User) PatchUserInfoById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	var req adaptersinboundhttpdto.PatchUserInfoByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := u.svcUser.SetUserInfoById(c.Request.Context(), id, req.Name, req.Bio); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PutUserPasswordMe godoc
// @Summary      Change my password
// @Description  Change the password of the currently authenticated user (requires old password)
// @Tags         Users
// @Accept       json
// @Security     BearerAuth
// @Param        body  body  adaptersinboundhttpdto.PutUserPasswordMeRequest  true  "Old and new password"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/me/password [put]
func (u *User) PutUserPasswordMe(c *gin.Context) {
	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	var req adaptersinboundhttpdto.PutUserPasswordMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := u.svcUser.SetUserPasswordById(c.Request.Context(), claims.Id, req.OldPassword, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserInvalidCredentials):
			c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Current password is incorrect"})
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PutUserPasswordById godoc
// @Summary      Override user password by ID
// @Description  Directly set a new password for a specific user (admin, no old password required)
// @Tags         Users
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  int                             true  "User ID"
// @Param        body  body  adaptersinboundhttpdto.PutUserPasswordByIdRequest  true  "New password"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/{id}/password [put]
func (u *User) PutUserPasswordById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	var req adaptersinboundhttpdto.PutUserPasswordByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := u.svcUser.SetUserPasswordDirectById(c.Request.Context(), id, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PutUserRoleById godoc
// @Summary      Set user role by ID
// @Description  Assign a role to a specific user
// @Tags         Users
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  int                        true  "User ID"
// @Param        body  body  adaptersinboundhttpdto.PutUserRoleByIdRequest  true  "New role"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/{id}/role [put]
func (u *User) PutUserRoleById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	var req adaptersinboundhttpdto.PutUserRoleByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := u.svcUser.SetUserRoleById(c.Request.Context(), id, req.Role); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteUserById godoc
// @Summary      Delete user by ID
// @Description  Permanently remove a user account
// @Tags         Users
// @Security     BearerAuth
// @Param        id  path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /users/{id} [delete]
func (u *User) DeleteUserById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	if err := u.svcUser.RemoveUserById(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
