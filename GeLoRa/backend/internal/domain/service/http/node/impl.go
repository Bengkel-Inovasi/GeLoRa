package domainservicehttpnode

import (
	"context"
	"errors"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
)

const path = "domain/service/http/node"

type impl struct {
	tx          *sqldt.Transactor
	log         portsoutboundlogging.Generic
	repoNode    portsoutboundrepository.Node
	repoSession portsoutboundrepository.Session
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoNode portsoutboundrepository.Node,
	repoSession portsoutboundrepository.Session,
) portsinboundhttp.Node {
	return &impl{
		tx:          tx,
		log:         log,
		repoNode:    repoNode,
		repoSession: repoSession,
	}
}

func (i *impl) GetNodeInfoById(ctx context.Context, id int64) (node *domainmodel.Node, err error) {
	const tag = path + "/GetNodeInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	node, err = i.repoNode.ReadNodeById(ctx, id)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read node by ID", logMeta(err))
		return nil, err
	}

	return node, nil
}

func (i *impl) GetNodesList(ctx context.Context, page int, limit int, cursorId *int64, search *string, isValidated *bool) (nodes []domainmodel.Node, total int, err error) {
	const tag = path + "/GetNodesList"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":        err.Error(),
			"page":         page,
			"limit":        limit,
			"cursor_id":    cursorId,
			"search":       search,
			"is_validated": isValidated,
		}
	}

	nodes, total, err = i.repoNode.ReadNodesByPagination(ctx, page, limit, cursorId, search, isValidated)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read nodes by pagination", logMeta(err))
		return nil, 0, err
	}

	return nodes, total, nil
}

func (i *impl) SetNodeInfoById(ctx context.Context, id int64, name *string, description *string) (err error) {
	const tag = path + "/SetNodeInfoById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":       err.Error(),
			"id":          id,
			"name":        name,
			"description": description,
		}
	}

	if err = i.repoNode.UpdateNodeInfoById(ctx, id, name, description); err != nil {
		i.log.Error(ctx, tag, "Failed to update node info", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "Node info updated", domainmodel.LogMeta{"id": id})
	return nil
}

func (i *impl) SetNodeValidate(ctx context.Context, id int64, validatedBy int64) (err error) {
	const tag = path + "/SetNodeValidate"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":        err.Error(),
			"id":           id,
			"validated_by": validatedBy,
		}
	}

	if err = i.repoNode.UpdateNodeValidation(ctx, id, true, validatedBy); err != nil {
		i.log.Error(ctx, tag, "Failed to update node validation", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "Node validated", domainmodel.LogMeta{"id": id, "validated_by": validatedBy})
	return nil
}

func (i *impl) RegisterClimber(ctx context.Context, mid string, name string, validatedBy int64) (nodeId int64, sessionId int64, err error) {
	const tag = path + "/RegisterClimber"

	node, err := i.repoNode.ReadNodeByMid(ctx, mid)
	if err != nil {
		if !errors.Is(err, domainmodel.ErrNodeNotFound) {
			i.log.Error(ctx, tag, "Failed to read node by MID", domainmodel.LogMeta{"error": err.Error(), "mid": mid})
			return 0, 0, err
		}
		nodeId, err = i.repoNode.CreateNode(ctx, mid, name)
		if err != nil {
			i.log.Error(ctx, tag, "Failed to create node", domainmodel.LogMeta{"error": err.Error(), "mid": mid})
			return 0, 0, err
		}
	} else {
		nodeId = node.Id
		if err = i.repoNode.UpdateNodeInfoById(ctx, nodeId, &name, nil); err != nil {
			i.log.Error(ctx, tag, "Failed to update node name", domainmodel.LogMeta{"error": err.Error(), "id": nodeId})
			return 0, 0, err
		}
	}

	if err = i.repoNode.UpdateNodeValidation(ctx, nodeId, true, validatedBy); err != nil {
		i.log.Error(ctx, tag, "Failed to validate node", domainmodel.LogMeta{"error": err.Error(), "id": nodeId})
		return 0, 0, err
	}

	// End any active session for this node before starting a new one
	_ = i.repoSession.UpdateSessionEndSession(ctx, nil, nil, &nodeId)

	sessionId, err = i.repoSession.CreateSession(ctx, validatedBy, nodeId)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to create session", domainmodel.LogMeta{"error": err.Error(), "node_id": nodeId})
		return 0, 0, err
	}

	i.log.Info(ctx, tag, "Climber registered", domainmodel.LogMeta{"node_id": nodeId, "session_id": sessionId, "mid": mid, "name": name})
	return nodeId, sessionId, nil
}

func (i *impl) RemoveNodeById(ctx context.Context, id int64) (err error) {
	const tag = path + "/RemoveNodeById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	if err = i.repoNode.DeleteNodeById(ctx, id); err != nil {
		i.log.Error(ctx, tag, "Failed to delete node by ID", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "Node removed", domainmodel.LogMeta{"id": id})
	return nil
}
