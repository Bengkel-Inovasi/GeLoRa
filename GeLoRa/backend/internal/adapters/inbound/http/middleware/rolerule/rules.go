package adaptersinboundhttpmiddlewarerolerule

import (
	"net/http"
	"strconv"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	middlewareauth "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/auth"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	"github.com/gin-gonic/gin"
)

func RuleAllowAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func RuleClaimsRoleHigherThanBodyRole(roleKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middlewareauth.GetJwtClaims(c)
		bodyRole, err := getRequestBodyRole(c, roleKey)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "Invalid request body",
				},
			)
			return
		}

		if claims.Role.Order() <= bodyRole.Order() {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "Insufficient role to assign this role",
				},
			)
			return
		}

		c.Next()
	}
}

func RuleClaimsRoleHigherThanUserIdPath(userIdParam string, userService portsinboundhttp.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middlewareauth.GetJwtClaims(c)
		userId, _ := strconv.ParseInt(c.Param(userIdParam), 10, 64)

		user, err := userService.GetUserInfoById(c.Request.Context(), userId)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "User not found",
				},
			)
			return
		}

		if claims.Role.Order() <= user.Role.Order() {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "Insufficient role",
				},
			)
			return
		}

		c.Next()
	}
}
