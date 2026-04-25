package adaptersinboundhttpmiddlewareauth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenFromBearerAuthorizationHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	authSplit := strings.Split(authHeader, " ")
	if len(authSplit) != 2 {
		return ""
	}
	if authSplit[0] != "Bearer" {
		return ""
	}

	return authSplit[1]
}
