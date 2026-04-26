package domainservicehttpsession

import (
	"context"
	"errors"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
)

const path = "domain/service/http/session"

type impl struct {
	tx          *sqldt.Transactor
	log         portsoutboundlogging.Generic
	repoSession portsoutboundrepository.Session
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoSession portsoutboundrepository.Session,
) portsinboundhttp.Session {
	return &impl{
		tx:          tx,
		log:         log,
		repoSession: repoSession,
	}
}

func (i *impl) AddSession(ctx context.Context, userId int64, nodeId int64) (id int64, err error) {
	const tag = path + "/AddSession"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":   err.Error(),
			"user_id": userId,
			"node_id": nodeId,
		}
	}

	userSession, err := i.repoSession.ReadSessionLatest(ctx, &userId, nil)
	if err != nil && !errors.Is(err, domainmodel.ErrSessionNotFound) {
		i.log.Error(ctx, tag, "Failed to read latest session for user", logMeta(err))
		return 0, err
	}
	if userSession != nil && userSession.EndedAt == nil {
		i.log.Error(ctx, tag, "User already has an active session", logMeta(domainmodel.ErrSessionAlreadyActive))
		return 0, domainmodel.ErrSessionAlreadyActive
	}

	nodeSession, err := i.repoSession.ReadSessionLatest(ctx, nil, &nodeId)
	if err != nil && !errors.Is(err, domainmodel.ErrSessionNotFound) {
		i.log.Error(ctx, tag, "Failed to read latest session for node", logMeta(err))
		return 0, err
	}
	if nodeSession != nil && nodeSession.EndedAt == nil {
		i.log.Error(ctx, tag, "Node already has an active session", logMeta(domainmodel.ErrSessionAlreadyActive))
		return 0, domainmodel.ErrSessionAlreadyActive
	}

	id, err = i.repoSession.CreateSession(ctx, userId, nodeId)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to create session", logMeta(err))
		return 0, err
	}

	i.log.Info(ctx, tag, "Session created", domainmodel.LogMeta{"id": id, "user_id": userId, "node_id": nodeId})
	return id, err
}

func (i *impl) GetSessionById(ctx context.Context, id int64) (session *domainmodel.Session, err error) {
	const tag = path + "/GetSessionById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	session, err = i.repoSession.ReadSessionById(ctx, id)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read session by ID", logMeta(err))
		return nil, err
	}

	return session, nil
}

func (i *impl) GetSessionsList(ctx context.Context, page int, limit int, cursorId *int64, userId *int64, nodeId *int64, active *bool) (sessions []domainmodel.Session, total int, err error) {
	const tag = path + "/GetSessionsList"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":     err.Error(),
			"page":      page,
			"limit":     limit,
			"cursor_id": cursorId,
			"user_id":   userId,
			"node_id":   nodeId,
			"active":    active,
		}
	}

	sessions, total, err = i.repoSession.ReadSessionsByPagination(ctx, page, limit, cursorId, userId, nodeId, active)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read sessions by pagination", logMeta(err))
		return nil, 0, err
	}

	return sessions, total, nil
}

func (i *impl) SetSessionEndSession(ctx context.Context, id *int64, userId *int64, nodeId *int64) (err error) {
	const tag = path + "/SetSessionEndSession"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":   err.Error(),
			"id":      id,
			"user_id": userId,
			"node_id": nodeId,
		}
	}

	if err = i.repoSession.UpdateSessionEndSession(ctx, id, userId, nodeId); err != nil {
		i.log.Error(ctx, tag, "Failed to end session", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "Session ended", domainmodel.LogMeta{"id": id, "user_id": userId, "node_id": nodeId})
	return nil
}

func (i *impl) RemoveSessionById(ctx context.Context, id int64) (err error) {
	const tag = path + "/RemoveSessionById"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error": err.Error(),
			"id":    id,
		}
	}

	if err = i.repoSession.DeleteSessionById(ctx, id); err != nil {
		i.log.Error(ctx, tag, "Failed to delete session by ID", logMeta(err))
		return err
	}

	i.log.Info(ctx, tag, "Session removed", domainmodel.LogMeta{"id": id})
	return nil
}
