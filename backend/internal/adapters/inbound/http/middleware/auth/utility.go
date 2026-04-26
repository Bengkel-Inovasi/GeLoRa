package adaptersinboundhttpmiddlewareauth

import (
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/gin-gonic/gin"
)

func GetJwtClaims(c *gin.Context) *domainmodel.UserAccessTokenClaims {
	claimsAny, exists := c.Get("claims")
	if !exists {
		return nil
	}
	claims, ok := claimsAny.(*domainmodel.UserAccessTokenClaims)
	if !ok {
		return nil
	}
	return claims
}
