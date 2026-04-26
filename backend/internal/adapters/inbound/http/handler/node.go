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

type Node struct {
	svcNode portsinboundhttp.Node
}

func NewNode(svcNode portsinboundhttp.Node) *Node {
	return &Node{
		svcNode: svcNode,
	}
}

// GetNodeInfoById godoc
// @Summary      Get node by ID
// @Description  Return the full details of a specific node
// @Tags         Nodes
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Node ID"
// @Success      200  {object}  adaptersinboundhttpdto.GetNodeInfoByIdResponse
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes/{id} [get]
func (n *Node) GetNodeInfoById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	node, err := n.svcNode.GetNodeInfoById(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Node not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertNodeModelToGetNodeInfoByIdResponse(node))
}

// GetNodesList godoc
// @Summary      List nodes
// @Description  Return a paginated list of nodes with optional filters
// @Tags         Nodes
// @Produce      json
// @Security     BearerAuth
// @Param        page          query     int     false  "Page number"             default(1)
// @Param        limit         query     int     false  "Items per page"          default(10)
// @Param        cursor_id     query     int     false  "Cursor ID for pagination"
// @Param        search        query     string  false  "Search by name or MID"
// @Param        is_validated  query     boolean false  "Filter by validation status"
// @Success      200           {object}  adaptersinboundhttpdto.GetNodesListResponse
// @Failure      400           {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401           {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500           {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes [get]
func (n *Node) GetNodesList(c *gin.Context) {
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

	var search *string
	if v := c.Query("search"); v != "" {
		search = &v
	}

	var isValidated *bool
	if v := c.Query("is_validated"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid is_validated"})
			return
		}
		isValidated = &parsed
	}

	nodes, total, err := n.svcNode.GetNodesList(c.Request.Context(), page, limit, cursorId, search, isValidated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		return
	}

	c.JSON(http.StatusOK, adaptersinboundhttpdto.ConvertNodeModelsToGetNodesListResponse(nodes, page, limit, total))
}

// PatchNodeInfoById godoc
// @Summary      Update node by ID
// @Description  Update the name and/or description of a specific node
// @Tags         Nodes
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  int                          true  "Node ID"
// @Param        body  body  adaptersinboundhttpdto.PatchNodeInfoByIdRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes/{id} [patch]
func (n *Node) PatchNodeInfoById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	var req adaptersinboundhttpdto.PatchNodeInfoByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	if err := n.svcNode.SetNodeInfoById(c.Request.Context(), id, req.Name, req.Description); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Node not found"})
		case errors.Is(err, domainmodel.ErrNodeNoChange):
			c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "No changes were provided"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PutNodeValidate godoc
// @Summary      Validate node by ID
// @Description  Mark a node as validated by the currently authenticated user
// @Tags         Nodes
// @Security     BearerAuth
// @Param        id  path  int  true  "Node ID"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes/{id}/validate [put]
func (n *Node) PutNodeValidate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	if err := n.svcNode.SetNodeValidate(c.Request.Context(), id, claims.Id); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Node not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PostRegisterClimber godoc
// @Summary      Register a climber to a LoRa device
// @Description  Find or create a node by MID, assign climber name, validate, and start a new session
// @Tags         Nodes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  adaptersinboundhttpdto.PostNodeRegisterRequest  true  "MID and climber name"
// @Success      201  {object}  adaptersinboundhttpdto.PostNodeRegisterResponse
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes/register [post]
func (n *Node) PostRegisterClimber(c *gin.Context) {
	var req adaptersinboundhttpdto.PostNodeRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: err.Error()})
		return
	}

	claims := middlewareauth.GetJwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, adaptersinboundhttpdto.CommonErrorResponse{Error: "Unauthorized"})
		return
	}

	nodeId, sessionId, err := n.svcNode.RegisterClimber(c.Request.Context(), req.Mid, req.Name, claims.Id)
	if err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{
				Error: "No device with LoRa address \"" + req.Mid + "\" has been detected. Make sure the wristband is powered on and transmitting.",
			})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.JSON(http.StatusCreated, adaptersinboundhttpdto.PostNodeRegisterResponse{
		NodeId:    nodeId,
		SessionId: sessionId,
	})
}

// DeleteNodeById godoc
// @Summary      Delete node by ID
// @Description  Permanently remove a node
// @Tags         Nodes
// @Security     BearerAuth
// @Param        id  path  int  true  "Node ID"
// @Success      204
// @Failure      400  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      401  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      404  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Failure      500  {object}  adaptersinboundhttpdto.CommonErrorResponse
// @Router       /nodes/{id} [delete]
func (n *Node) DeleteNodeById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, adaptersinboundhttpdto.CommonErrorResponse{Error: "Invalid id"})
		return
	}

	if err := n.svcNode.RemoveNodeById(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, domainmodel.ErrNodeNotFound):
			c.JSON(http.StatusNotFound, adaptersinboundhttpdto.CommonErrorResponse{Error: "Node not found"})
		default:
			c.JSON(http.StatusInternalServerError, adaptersinboundhttpdto.CommonErrorResponse{Error: "Something wrong happened"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
