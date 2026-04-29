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

type Alert struct {
	svcAlert portsinboundhttp.Alert
}

func NewAlert(svcAlert portsinboundhttp.Alert) *Alert {
	return &Alert{svcAlert: svcAlert}
}

// PostAlert creates an emergency alert from any authenticated user.
func (a *Alert) PostAlert(c *gin.Context) {
	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	var req adaptersinboundhttpdto.PostAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	id, err := a.svcAlert.AddAlert(c.Request.Context(), claims.Id, req.NodeId, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusCreated, adaptersinboundhttpdto.PostAlertResponse{Id: id})
}

// GetAlertsList returns alerts (unacknowledged by default) for operators.
func (a *Alert) GetAlertsList(c *gin.Context) {
	unacknowledgedOnly := c.DefaultQuery("unacknowledged", "true") == "true"

	alerts, err := a.svcAlert.GetAlerts(c.Request.Context(), unacknowledgedOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertAlertModelsToGetAlertsListResponse(alerts))
}

// PutAlertAcknowledge marks an alert as acknowledged/resolved.
func (a *Alert) PutAlertAcknowledge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	if err := a.svcAlert.AcknowledgeAlert(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrAlertNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Alert not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
