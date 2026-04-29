package domainservicehttpalert

import (
	"context"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
)

const path = "domain/service/http/alert"

type impl struct {
	tx        *sqldt.Transactor
	log       portsoutboundlogging.Generic
	repoAlert portsoutboundrepository.Alert
}

func NewImpl(
	tx *sqldt.Transactor,
	log portsoutboundlogging.Generic,
	repoAlert portsoutboundrepository.Alert,
) portsinboundhttp.Alert {
	return &impl{tx: tx, log: log, repoAlert: repoAlert}
}

func (i *impl) AddAlert(ctx context.Context, userId int64, nodeId *int64, message string) (id int64, err error) {
	const tag = path + "/AddAlert"
	id, err = i.repoAlert.CreateAlert(ctx, userId, nodeId, message)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to create alert", domainmodel.LogMeta{"error": err.Error(), "user_id": userId})
		return 0, err
	}
	i.log.Info(ctx, tag, "Alert created", domainmodel.LogMeta{"id": id, "user_id": userId})
	return id, nil
}

func (i *impl) GetAlerts(ctx context.Context, unacknowledgedOnly bool) (alerts []domainmodel.Alert, err error) {
	const tag = path + "/GetAlerts"
	alerts, err = i.repoAlert.ReadAlerts(ctx, unacknowledgedOnly)
	if err != nil {
		i.log.Error(ctx, tag, "Failed to read alerts", domainmodel.LogMeta{"error": err.Error()})
		return nil, err
	}
	return alerts, nil
}

func (i *impl) AcknowledgeAlert(ctx context.Context, id int64) (err error) {
	const tag = path + "/AcknowledgeAlert"
	if err = i.repoAlert.UpdateAlertAcknowledge(ctx, id); err != nil {
		i.log.Error(ctx, tag, "Failed to acknowledge alert", domainmodel.LogMeta{"error": err.Error(), "id": id})
		return err
	}
	i.log.Info(ctx, tag, "Alert acknowledged", domainmodel.LogMeta{"id": id})
	return nil
}
