package domainappseeder

import (
	"context"

	databaseseeder "github.com/Bengkel-Inovasi/GeLoRa/backend/database/seeder"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/pwd"
)

const path = "domain/app/seeder"

type Seeder struct {
	infrastructure *Infrastructure
	wiring         *Wiring
}

func Run() int {
	const tag = path + "/Run"
	s := &Seeder{}
	ctx := context.Background()

	if err := config.LoadEnv(); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to load .env", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := s.NewInfrastructure(ctx); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create infrastructure", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := s.NewWiring(ctx); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create wiring", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := s.seed(ctx); err != nil {
		s.wiring.log.Error(ctx, tag, "Seeding failed", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	s.wiring.log.Info(ctx, tag, "Seeding completed successfully", nil)
	return 0
}

func (s *Seeder) seed(ctx context.Context) error {
	const tag = path + "/seed"

	seedFile, err := databaseseeder.LoadUserSeedFile()
	if err != nil {
		s.wiring.log.Error(ctx, tag, "Failed to load user seed file", domainmodel.LogMeta{"error": err.Error()})
		return err
	}

	for _, u := range seedFile.Seed {
		passwordHash, err := pwd.Hash(u.Password)
		if err != nil {
			s.wiring.log.Error(ctx, tag, "Failed to hash password", domainmodel.LogMeta{"error": err.Error(), "username": u.Username})
			return err
		}

		if err = s.wiring.repoUser.CreateUserSeed(
			ctx,
			u.Id,
			u.Name,
			u.Username,
			passwordHash,
			domainmodel.UserRole(u.Role),
			u.Bio,
		); err != nil {
			s.wiring.log.Error(ctx, tag, "Failed to seed user", domainmodel.LogMeta{"error": err.Error(), "username": u.Username})
			return err
		}

		s.wiring.log.Info(ctx, tag, "User seeded", domainmodel.LogMeta{"id": u.Id, "username": u.Username})
	}

	return nil
}
