package adaptersinboundhttpdto

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

// -- POST /alerts ------------------------------------------------------------

type PostAlertRequest struct {
	NodeId  *int64 `json:"node_id"  example:"3"`
	Message string `json:"message"  example:"My friend needs help!"`
}

type PostAlertResponse struct {
	Id int64 `json:"id" example:"1"`
}

// -- GET /alerts -------------------------------------------------------------

type GetAlertsListData struct {
	Id             int64      `json:"id"`
	UserId         *int64     `json:"user_id"`
	NodeId         *int64     `json:"node_id"`
	Message        string     `json:"message"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type GetAlertsListResponse struct {
	Data []GetAlertsListData `json:"data"`
}

func ConvertAlertModelsToGetAlertsListResponse(alerts []domainmodel.Alert) *GetAlertsListResponse {
	data := make([]GetAlertsListData, len(alerts))
	for i, a := range alerts {
		data[i] = GetAlertsListData{
			Id:             a.Id,
			UserId:         a.UserId,
			NodeId:         a.NodeId,
			Message:        a.Message,
			AcknowledgedAt: a.AcknowledgedAt,
			CreatedAt:      a.CreatedAt,
		}
	}
	return &GetAlertsListResponse{Data: data}
}
