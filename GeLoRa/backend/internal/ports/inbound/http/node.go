package portsinboundhttp

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Node interface {
	GetNodeInfoById(ctx context.Context, id int64) (node *domainmodel.Node, err error)
	GetNodesList(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (nodes []domainmodel.Node, total int, err error)
	SetNodeInfoById(ctx context.Context, id int64, name *string, description *string) (err error)
	SetNodeValidate(ctx context.Context, id int64, validatedBy int64) (err error)
	RemoveNodeById(ctx context.Context, id int64) (err error)
	RegisterClimber(ctx context.Context, mid string, name string, validatedBy int64) (nodeId int64, sessionId int64, err error)
}
