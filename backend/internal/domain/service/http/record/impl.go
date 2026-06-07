package domainservicehttprecord

import (
	"context"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
)

const path = "domain/service/http/record"

type impl struct {
	tx         *sqldt.Transactor
	log        portsoutboundlogging.Generic
	repoRecord portsoutboundrepository.Record
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoRecord portsoutboundrepository.Record,
) portsinboundhttp.Record {
	return &impl{
		tx:         tx,
		log:        log,
		repoRecord: repoRecord,
	}
}

func (i *impl) GetRecords(ctx context.Context, sessionId *int64, userId *int64, nodeId *int64, startTime *time.Time, endTime *time.Time) (records []domainmodel.Record, err error) {
	const tag = path + "/GetRecords"
	logMeta := func(err error) domainmodel.LogMeta {
		return domainmodel.LogMeta{
			"error":      err.Error(),
			"session_id": sessionId,
			"user_id":    userId,
			"node_id":    nodeId,
			"start_time": startTime,
			"end_time":   endTime,
		}
	}

	records, err = i.repoRecord.ReadRecords(ctx, sessionId, userId, nodeId, startTime, endTime)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read records", logMeta(err))
		return nil, err
	}

	i.log.Info(ctx, tag, "Records retrieved", domainmodel.LogMeta{
		"count":      len(records),
		"node_id":    nodeId,
		"start_time": startTime,
	})

	return records, nil
}
