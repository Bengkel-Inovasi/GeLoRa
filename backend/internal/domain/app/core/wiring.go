package domainappcore

import (
	"context"

	adaptersoutboundlogginggeneric "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/logging/generic"
	adaptersoutboundrepositoryalert "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/alert"
	adaptersoutboundrepositorynode "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/node"
	adaptersoutboundrepositoryrecord "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/record"
	adaptersoutboundrepositorysession "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/session"
	adaptersoutboundrepositoryuser "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/user"
	domainservicehttpalert "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/http/alert"
	domainservicehttpnode "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/http/node"
	domainservicehttprecord "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/http/record"
	domainservicehttpsession "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/http/session"
	domainservicehttpuser "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/http/user"
	domainservicemqttrecord "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/service/mqtt/record"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsinboundmqtt "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/mqtt"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
)

type Wiring struct {
	log             portsoutboundlogging.Generic
	repoUser        portsoutboundrepository.User
	repoNode        portsoutboundrepository.Node
	repoSession     portsoutboundrepository.Session
	repoRecord      portsoutboundrepository.Record
	repoAlert       portsoutboundrepository.Alert
	svcHttpUser     portsinboundhttp.User
	svcHttpNode     portsinboundhttp.Node
	svcHttpSession  portsinboundhttp.Session
	svcHttpRecord   portsinboundhttp.Record
	svcHttpAlert    portsinboundhttp.Alert
	svcMqttRecord   portsinboundmqtt.Record
}

func (c *Core) NewWiring(ctx context.Context) (err error) {
	const tag = path + "/NewWiring"

	// Logging
	log := adaptersoutboundlogginggeneric.NewSlogImpl(c.infrastructure.log)

	// Repositories
	repoUser := adaptersoutboundrepositoryuser.NewSqliteImpl(
		c.infrastructure.sqlDT,
		c.infrastructure.sqrQ,
		c.infrastructure.sqrD,
		log,
	)
	repoNode := adaptersoutboundrepositorynode.NewSqliteImpl(
		c.infrastructure.sqlDT,
		c.infrastructure.sqrQ,
		c.infrastructure.sqrD,
		log,
	)
	repoSession := adaptersoutboundrepositorysession.NewSqliteImpl(
		c.infrastructure.sqlDT,
		c.infrastructure.sqrQ,
		c.infrastructure.sqrD,
		log,
	)
	repoRecord := adaptersoutboundrepositoryrecord.NewSqliteImpl(
		c.infrastructure.sqlDT,
		c.infrastructure.sqrQ,
		c.infrastructure.sqrD,
		log,
	)
	repoAlert := adaptersoutboundrepositoryalert.NewSqliteImpl(
		c.infrastructure.sqlDT,
		c.infrastructure.sqrQ,
		c.infrastructure.sqrD,
		log,
	)
	log.Info(ctx, tag, "Repositories initiated", nil)

	// Services
	svcHttpUser := domainservicehttpuser.NewImpl(c.infrastructure.sqlTX, log, repoUser)
	svcHttpNode := domainservicehttpnode.NewImpl(c.infrastructure.sqlTX, log, repoNode, repoSession)
	svcHttpSession := domainservicehttpsession.NewImpl(c.infrastructure.sqlTX, log, repoSession)
	svcHttpRecord := domainservicehttprecord.NewImpl(c.infrastructure.sqlTX, log, repoRecord)
	svcHttpAlert := domainservicehttpalert.NewImpl(c.infrastructure.sqlTX, log, repoAlert)
	svcMqttRecord := domainservicemqttrecord.NewImpl(c.infrastructure.sqlTX, log, repoRecord, repoNode, repoSession)
	log.Info(ctx, tag, "Services initiated", nil)

	// Inject
	c.wiring = &Wiring{
		log:            log,
		repoUser:       repoUser,
		repoNode:       repoNode,
		repoSession:    repoSession,
		repoRecord:     repoRecord,
		repoAlert:      repoAlert,
		svcHttpUser:    svcHttpUser,
		svcHttpNode:    svcHttpNode,
		svcHttpSession: svcHttpSession,
		svcHttpRecord:  svcHttpRecord,
		svcHttpAlert:   svcHttpAlert,
		svcMqttRecord:  svcMqttRecord,
	}
	log.Info(ctx, tag, "Wiring created successfully", nil)
	return nil
}
