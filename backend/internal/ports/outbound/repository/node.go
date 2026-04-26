package portsoutboundrepository

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

type Node interface {
	CreateNode(ctx context.Context, mid string, name string) (id int64, err error)
	ReadNodeById(ctx context.Context, id int64) (node *domainmodel.Node, err error)
	ReadNodeByMid(ctx context.Context, mid string) (node *domainmodel.Node, err error)
	ReadNodesByPagination(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (nodes []domainmodel.Node, total int, err error)
	UpdateNodeInfoById(ctx context.Context, id int64, name *string, description *string) (err error)
	UpdateNodeValidation(ctx context.Context, id int64, validation bool, validatedBy int64) (err error)
	DeleteNodeById(ctx context.Context, id int64) (err error)
}
