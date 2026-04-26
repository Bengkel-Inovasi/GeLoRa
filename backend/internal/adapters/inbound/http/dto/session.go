package adaptersinboundhttpdto

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

// -- POST /sessions ----------------------------------------------------------

type PostSessionRequest struct {
	UserId int64 `json:"user_id" example:"1"`
	NodeId int64 `json:"node_id" example:"3"`
}

type PostSessionResponse struct {
	Id int64 `json:"id" example:"42"`
}

// -- GET /sessions/:id -------------------------------------------------------

type GetSessionByIdResponse struct {
	Id        int64      `json:"id"         example:"42"`
	UserId    *int64     `json:"user_id"    example:"1"`
	NodeId    *int64     `json:"node_id"    example:"3"`
	StartedAt time.Time  `json:"started_at" example:"2024-01-01T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at"   example:"2024-01-01T01:00:00Z"`
}

func ConvertSessionModelToGetSessionByIdResponse(session *domainmodel.Session) *GetSessionByIdResponse {
	return &GetSessionByIdResponse{
		Id:        session.Id,
		UserId:    session.UserId,
		NodeId:    session.NodeId,
		StartedAt: session.StartedAt,
		EndedAt:   session.EndedAt,
	}
}

// -- GET /sessions -----------------------------------------------------------

type GetSessionsListData struct {
	Id        int64      `json:"id"         example:"42"`
	UserId    *int64     `json:"user_id"    example:"1"`
	NodeId    *int64     `json:"node_id"    example:"3"`
	StartedAt time.Time  `json:"started_at" example:"2024-01-01T00:00:00Z"`
	EndedAt   *time.Time `json:"ended_at"   example:"2024-01-01T01:00:00Z"`
}

type GetSessionsListResponse struct {
	Data       []GetSessionsListData `json:"data"`
	Pagination CommonPaginationMeta  `json:"pagination"`
}

func ConvertSessionModelsToGetSessionsListResponse(sessions []domainmodel.Session, page, limit, total int) *GetSessionsListResponse {
	data := make([]GetSessionsListData, len(sessions))
	for i, s := range sessions {
		data[i] = GetSessionsListData{
			Id:        s.Id,
			UserId:    s.UserId,
			NodeId:    s.NodeId,
			StartedAt: s.StartedAt,
			EndedAt:   s.EndedAt,
		}
	}
	return &GetSessionsListResponse{
		Data: data,
		Pagination: CommonPaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

// -- PUT /sessions/end -------------------------------------------------------

type PutSessionEndRequest struct {
	Id     *int64 `json:"id"      example:"42"`
	UserId *int64 `json:"user_id" example:"1"`
	NodeId *int64 `json:"node_id" example:"3"`
}
