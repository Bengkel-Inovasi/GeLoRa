package adaptersinboundhttphandler

import (
	"errors"
	"net/http"
	"strconv"

	adaptersinboundhttpdto "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/dto"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	"github.com/gin-gonic/gin"
)

type Session struct {
	svcSession portsinboundhttp.Session
}

func NewSession(svcSession portsinboundhttp.Session) *Session {
	return &Session{
		svcSession: svcSession,
	}
}

// PostSession godoc
// @Summary      Start a session
// @Description  Create a new active session for a user on a node
// @Tags         Sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      adaptersinboundhttpdto.PostSessionRequest  true  "Session details"
// @Success      201   {object}  adaptersinboundhttpdto.PostSessionResponse
// @Failure      400   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      409   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500   {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /sessions [post]
func (s *Session) PostSession(c *gin.Context) {
	var req adaptersinboundhttpdto.PostSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	id, err := s.svcSession.AddSession(c.Request.Context(), req.UserId, req.NodeId)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrSessionAlreadyActive):
			c.JSON(http.StatusConflict, adaptersinboundhttpdto.CommonErrorResponse{Error: "A session is already active for this user and node"})
		case errors.Is(err, domainmodel.ErrUserNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "User not found"})
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Node not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusCreated, adaptersinboundhttpdto.PostSessionResponse{Id: id})
}

// GetSessionById godoc
// @Summary      Get session by ID
// @Description  Return the details of a specific session
// @Tags         Sessions
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Session ID"
// @Success      200  {object}  adaptersinboundhttpdto.GetSessionByIdResponse
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /sessions/{id} [get]
func (s *Session) GetSessionById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	session, err := s.svcSession.GetSessionById(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Session not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertSessionModelToGetSessionByIdResponse(session))
}

// GetSessionsList godoc
// @Summary      List sessions
// @Description  Return a paginated list of sessions with optional filters
// @Tags         Sessions
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number"             default(1)
// @Param        limit      query     int     false  "Items per page"          default(10)
// @Param        cursor_id  query     int     false  "Cursor ID for pagination"
// @Param        user_id    query     int     false  "Filter by user ID"
// @Param        node_id    query     int     false  "Filter by node ID"
// @Param        active     query     boolean false  "Filter by active status"
// @Success      200        {object}  adaptersinboundhttpdto.GetSessionsListResponse
// @Failure      400        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500        {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /sessions [get]
func (s *Session) GetSessionsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var cursorId *int64
	if v := c.Query("cursor_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid cursor_id"})
			return
		}
		cursorId = &parsed
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

	var active *bool
	if v := c.Query("active"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid active"})
			return
		}
		active = &parsed
	}

	sessions, total, err := s.svcSession.GetSessionsList(c.Request.Context(), page, limit, cursorId, userId, nodeId, active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertSessionModelsToGetSessionsListResponse(sessions, page, limit, total))
}

// PutSessionEnd godoc
// @Summary      End a session
// @Description  Terminate a session identified by session ID, user ID, or node ID (at least one required)
// @Tags         Sessions
// @Accept       json
// @Security     BearerAuth
// @Param        body  body  adaptersinboundhttpdto.PutSessionEndRequest  true  "Session identifier(s)"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /sessions/end [put]
func (s *Session) PutSessionEnd(c *gin.Context) {
	var req adaptersinboundhttpdto.PutSessionEndRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := s.svcSession.SetSessionEndSession(c.Request.Context(), req.Id, req.UserId, req.NodeId); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Session not found"})
		case errors.Is(err, domainmodel.ErrSessionNoIdentifier):
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "At least one session identifier is required"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteSessionById godoc
// @Summary      Delete session by ID
// @Description  Permanently remove a session record
// @Tags         Sessions
// @Security     BearerAuth
// @Param        id  path  int  true  "Session ID"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /sessions/{id} [delete]
func (s *Session) DeleteSessionById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	if err := s.svcSession.RemoveSessionById(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Session not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
