package domainappseeder

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Masterminds/squirrel"
	_ "modernc.org/sqlite"
)

type Infrastructure struct {
	log   *slog.Logger
	sqlDT sqldt.Sqldt
	sqrD  *squirrel.StatementBuilderType
	sqrQ  *squirrel.StatementBuilderType
}

func (s *Seeder) NewInfrastructure(ctx context.Context) error {
	const tag = path + "/NewInfrastructure"

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	sqliteDB, err := sql.Open("sqlite", "file:"+config.SQLiteDatabase+"?mode=rw&cache=shared&_journal_mode=WAL")
	if err != nil {
		logTmp(log).Error(ctx, tag, "Failed to open SQLite database", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	if err = sqliteDB.PingContext(ctx); err != nil {
		logTmp(log).Error(ctx, tag, "SQLite database not reachable, ensure the file exists and migrations have been run", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	logTmp(log).Info(ctx, tag, "SQLite Database connected", domainmodel.LogMeta{"path": config.SQLiteDatabase})

	sqrD := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqrQ := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	s.infrastructure = &Infrastructure{
		log:   log,
		sqlDT: sqldt.NewSqldt(sqliteDB),
		sqrD:  &sqrD,
		sqrQ:  &sqrQ,
	}
	logTmp(log).Info(ctx, tag, "Infrastructures created successfully", nil)
	return nil
}
