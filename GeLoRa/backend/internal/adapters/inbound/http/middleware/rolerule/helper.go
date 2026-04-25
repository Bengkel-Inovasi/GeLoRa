package adaptersinboundhttpmiddlewarerolerule

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func getRequestBodyRole(c *gin.Context, roleKey string) (domainmodel.UserRole, error) {
	var bodyMap map[string]any

	if c.Request.Body == nil {
		return domainmodel.UserRoleClient, errors.New("no body")
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return domainmodel.UserRoleClient, err
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		return domainmodel.UserRoleClient, err
	}

	roleVal, exists := bodyMap[roleKey]
	if !exists {
		return domainmodel.UserRoleClient, errors.New("no role in the body")
	}

	roleStr, ok := roleVal.(string)
	if !ok {
		return domainmodel.UserRoleClient, errors.New("role is not a string")
	}

	if err := utils.ValidateRole(roleStr); err != nil {
		return domainmodel.UserRoleClient, err
	}

	return domainmodel.UserRole(roleStr), nil
}
