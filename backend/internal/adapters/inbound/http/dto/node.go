package adaptersinboundhttpdto

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/utils"
)

// -- GET /nodes/:id ----------------------------------------------------------

type GetNodeInfoByIdResponse struct {
	Id          int64      `json:"id"           example:"1"`
	Mid         string     `json:"mid"          example:"node-abc123"`
	Name        string     `json:"name"         example:"Main Gateway"`
	Description string     `json:"description"  example:"Primary entry node."`
	IsValidated bool       `json:"is_validated" example:"true"`
	ValidatedBy *int64     `json:"validated_by" example:"2"`
	CreatedAt   time.Time  `json:"created_at"   example:"2024-01-01T00:00:00Z"`
	UpdatedAt   *time.Time `json:"updated_at"   example:"2024-06-01T00:00:00Z"`
}

func ConvertNodeModelToGetNodeInfoByIdResponse(node *domainmodel.Node) *GetNodeInfoByIdResponse {
	return &GetNodeInfoByIdResponse{
		Id:          node.Id,
		Mid:         node.Mid,
		Name:        node.Name,
		Description: node.Description,
		IsValidated: node.IsValidated,
		ValidatedBy: node.ValidatedBy,
		CreatedAt:   node.CreatedAt,
		UpdatedAt:   node.UpdatedAt,
	}
}

// -- GET /nodes --------------------------------------------------------------

type GetNodesListData struct {
	Id          int64  `json:"id"           example:"1"`
	Mid         string `json:"mid"          example:"node-abc123"`
	Name        string `json:"name"         example:"Main Gateway"`
	IsValidated bool   `json:"is_validated" example:"true"`
}

type GetNodesListResponse struct {
	Data       []GetNodesListData   `json:"data"`
	Pagination CommonPaginationMeta `json:"pagination"`
}

func ConvertNodeModelsToGetNodesListResponse(nodes []domainmodel.Node, page, limit, total int) *GetNodesListResponse {
	data := make([]GetNodesListData, len(nodes))
	for i, n := range nodes {
		data[i] = GetNodesListData{
			Id:          n.Id,
			Mid:         n.Mid,
			Name:        n.Name,
			IsValidated: n.IsValidated,
		}
	}
	return &GetNodesListResponse{
		Data: data,
		Pagination: CommonPaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

// -- PATCH /nodes/:id --------------------------------------------------------

type PatchNodeInfoByIdRequest struct {
	Name        *string `json:"name"        example:"Main Gateway"`
	Description *string `json:"description" example:"Primary entry node."`
}

func (r *PatchNodeInfoByIdRequest) Validate() error {
	if r.Name != nil {
		if err := utils.ValidateName(*r.Name); err != nil {
			return err
		}
	}
	return nil
}

// -- PUT /nodes/:id/validate -------------------------------------------------

type PutNodeValidateRequest struct {
	ValidatedBy int64 `json:"validated_by" example:"2"`
}

// -- POST /nodes/register ----------------------------------------------------

type PostNodeRegisterRequest struct {
	Mid  string `json:"mid"  binding:"required" example:"PENDAKI_01"`
	Name string `json:"name" binding:"required" example:"Ahmad Fauzi"`
}

type PostNodeRegisterResponse struct {
	NodeId    int64 `json:"node_id"    example:"1"`
	SessionId int64 `json:"session_id" example:"5"`
}
