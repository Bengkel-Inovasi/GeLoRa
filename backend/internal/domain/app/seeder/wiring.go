package domainappseeder

import (
	"context"

	adaptersoutboundlogginggeneric "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/logging/generic"
	adaptersoutboundrepositoryuser "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/outbound/repository/user"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	portsoutboundrepository "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/repository"
)

type Wiring struct {
	log      portsoutboundlogging.Generic
	repoUser portsoutboundrepository.User
}

func (s *Seeder) NewWiring(ctx context.Context) error {
	const tag = path + "/NewWiring"

	log := adaptersoutboundlogginggeneric.NewSlogImpl(s.infrastructure.log)

	repoUser := adaptersoutboundrepositoryuser.NewSqliteImpl(
		s.infrastructure.sqlDT,
		s.infrastructure.sqrQ,
		s.infrastructure.sqrD,
		log,
	)

	s.wiring = &Wiring{
		log:      log,
		repoUser: repoUser,
	}
	log.Info(ctx, tag, "Wiring created successfully", nil)
	return nil
}
