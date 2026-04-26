package adaptersinboundhttphandler

import (
	"net/http"
	"strconv"
	"time"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	"github.com/gin-gonic/gin"
)

type Record struct {
	svcRecord portsinboundhttp.Record
}

func NewRecord(svcRecord portsinboundhttp.Record) *Record {
	return &Record{
		svcRecord: svcRecord,
	}
}

// GetRecords godoc
// @Summary      List records
// @Description  Return telemetry records with optional filters by session, user, node, and time range
// @Tags         Records
// @Produce      json
// @Security     BearerAuth
// @Param        session_id  query     int     false  "Filter by session ID"
// @Param        user_id     query     int     false  "Filter by user ID"
// @Param        node_id     query     int     false  "Filter by node ID"
// @Param        start_time  query     string  false  "Start of time range (RFC3339, e.g. 2024-01-01T00:00:00Z)"
// @Param        end_time    query     string  false  "End of time range (RFC3339, e.g. 2024-01-01T23:59:59Z)"
// @Success      200         {object}  adaptersinboundhttpdto.GetRecordsResponse
// @Failure      400         {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401         {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500         {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /records [get]
func (r *Record) GetRecords(c *gin.Context) {
	var sessionId *int64
	if v := c.Query("session_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid session_id"})
			return
		}
		sessionId = &parsed
	}

	var userId *int64
	if v := c.Query("user_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid user_id"})
			return
		}
		userId = &parsed
	}

	var nodeId *int64
	if v := c.Query("node_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid node_id"})
			return
		}
		nodeId = &parsed
	}

	var startTime *time.Time
	if v := c.Query("start_time"); v != "" {
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid start_time, use RFC3339 format"})
			return
		}
		startTime = &parsed
	}

	var endTime *time.Time
	if v := c.Query("end_time"); v != "" {
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid end_time, use RFC3339 format"})
			return
		}
		endTime = &parsed
	}

	records, err := r.svcRecord.GetRecords(c.Request.Context(), sessionId, userId, nodeId, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertRecordModelsToGetRecordsResponse(records))
}
