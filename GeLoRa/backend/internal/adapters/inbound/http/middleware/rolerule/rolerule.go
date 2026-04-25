package adaptersinboundhttpmiddlewarerolerule

import (
	"net/http"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	middlewareauth "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/auth"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type RoleRuler struct {
	Role domainmodel.UserRole
	Rule gin.HandlerFunc
}

func RoleRule(roleRules []RoleRuler) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middlewareauth.GetJwtClaims(c)
		if claims == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "Unauthorized",
				},
			)
			return
		}

		for _, roleRule := range roleRules {
			if roleRule.Role == claims.Role {
				roleRule.Rule(c)
				return
			}
		}

		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			adaptersinboundhttpdto.CommonErrorResponse{
				Error: "Unauthorized",
			},
		)
	}
}
