package adaptersinboundhttpmiddlewareauth

import (
	"errors"
	"net/http"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/token"
	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := getTokenFromBearerAuthorizationHeader(c)
		if tkn == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				adaptersinboundhttpdto.CommonErrorResponse{
					Error: "Valid bearer authorization header is required",
				},
			)
			return
		}

		claims := &domainmodel.UserAccessTokenClaims{}
		if err := token.ValidateAccess(tkn, claims); err != nil {
			if errors.Is(err, token.ErrInvalid) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					adaptersinboundhttpdto.CommonErrorResponse{
						Error: "Invalid access token",
					},
				)
				return
			}
			if errors.Is(err, token.ErrExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					adaptersinboundhttpdto.CommonErrorResponse{
						Error: "Expired access token",
					},
				)
				return
			}
		}

		c.Set("claims", claims)
		c.Next()
	}
}
