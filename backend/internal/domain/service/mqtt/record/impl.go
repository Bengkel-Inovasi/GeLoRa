package domainservicemqttrecord

import (
	"context"
	"errors"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundmqtt "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/mqtt"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
)

const path = "domain/service/mqtt/record"

type impl struct {
	tx          *sqldt.Transactor
	log         portsoutboundlogging.Generic
	repoRecord  portsoutboundrepository.Record
	repoNode    portsoutboundrepository.Node
	repoSession portsoutboundrepository.Session
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoRecord portsoutboundrepository.Record,
	repoNode portsoutboundrepository.Node,
	repoSession portsoutboundrepository.Session,
) portsinboundmqtt.Record {
	return &impl{
		tx:          tx,
		log:         log,
		repoRecord:  repoRecord,
		repoNode:    repoNode,
		repoSession: repoSession,
	}
}

func (i *impl) AddRecord(ctx context.Context, mid string, time time.Time, heartRate *float64, bodyTemperature *float64, latitude *float64, longitude *float64) (id int64, err error) {
	const tag = path + "/AddRecord"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":            err.Error(),
			"mid":              mid,
			"time":             time,
			"heart_rate":       heartRate,
			"body_temperature": bodyTemperature,
			"latitude":         latitude,
			"longitude":        longitude,
		}
	}

	// 1. Resolve the node and its validation from MID.
	node, err := i.repoNode.ReadNodeByMid(ctx, mid)
	if err != nil {
		if !errors.Is(err, domainmodel.ErrNodeNotFound) {
			i.log.Error(ctx, tag, "Failed to read node by MID", logMeta(err))
			return 0, err
		}

		nodeId, err := i.repoNode.CreateNode(ctx, mid, mid)
		if err != nil && !errors.Is(err, domainmodel.ErrNodeMidExists) {
			i.log.Error(ctx, tag, "Failed to create node", logMeta(err))
			return 0, err
		}

		i.log.Info(ctx, tag, "New node registered", domainmodel.LogMeta{
			"node_id": nodeId,
			"mid":     mid,
		})
		return 0, nil
	}
	if !node.IsValidated {
		return 0, nil
	}

	// 2. Check for an active session and save the record atomically.
	session, err := i.repoSession.ReadSessionLatest(ctx, nil, &node.Id)
	if err != nil {
		if errors.Is(err, domainmodel.ErrSessionNotFound) {
			return 0, nil
		}
		i.log.Error(ctx, tag, "Failed to read latest session for node", logMeta(err))
		return 0, err
	}
	if session.EndedAt != nil {
		return 0, nil
	}

	id, err = i.repoRecord.CreateRecord(ctx, session.Id, time, heartRate, bodyTemperature, latitude, longitude)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to create record", logMeta(err))
		return 0, err
	}

	i.log.Info(ctx, tag, "Record created", domainmodel.LogMeta{"id": id, "session_id": session.Id, "mid": mid})
	return id, err
}
