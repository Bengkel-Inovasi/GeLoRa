package domainappcore

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/sqldt"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/pkg/token"
	"github.com/Masterminds/squirrel"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Infrastructure struct {
	log        *slog.Logger
	sqlDT      sqldt.Sqldt
	sqlTX      *sqldt.Transactor
	sqrD       *squirrel.StatementBuilderType
	sqrQ       *squirrel.StatementBuilderType
	mqttClient mqtt.Client
	ginEngine  *gin.Engine
}

func (c *Core) NewInfrastructure(ctx context.Context) (err error) {
	const tag = path + "/NewInfrastructure"

	// Log
	log := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: getSlogLevelFromDomainLogLevel(config.AppLogLevel)},
	))

	// Token config
	token.Config(
		config.TokenIssuer,
		config.TokenAccessPrivateKey,
		config.TokenRefreshPrivateKey,
		config.TokenAccessTTL,
		config.TokenRefreshTTL,
	)
	logTmp(log).Info(ctx, tag, "Token config set", nil)

	// SQLite Database
	sqliteDB, err := sql.Open("sqlite", getSQLiteConnectionPath(config.SQLiteDatabase))
	if err != nil {
		logTmp(log).Error(ctx, tag, "Failed to open SQLite Database", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	if err = sqliteDB.PingContext(ctx); err != nil {
		logTmp(log).Error(ctx, tag, "SQLite Database unreachable, make sure the database file exists", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	sqlDT := sqldt.NewSqldt(sqliteDB)
	sqlTX := sqldt.NewTransactor(sqliteDB)
	sqrD := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqrQ := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	logTmp(log).Info(ctx, tag, "SQLite Database connected", domainmodel.LogMeta{"path": config.SQLiteDatabase})

	// MQTT Client
	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(config.MQTTBrokerURL)
	mqttOpts.SetClientID(config.MQTTBrokerClientID)
	mqttOpts.SetUsername(config.MQTTBrokerUsername)
	mqttOpts.SetPassword(config.MQTTBrokerPassword)
	mqttOpts.SetAutoReconnect(true)
	mqttOpts.SetMaxReconnectInterval(5 * time.Minute)
	mqttOpts.SetConnectTimeout(10 * time.Second)
	mqttOpts.SetKeepAlive(60 * time.Second)
	mqttOpts.SetCleanSession(false)
	mqttOpts.SetWriteTimeout(10 * time.Second)
	mqttOpts.OnConnect = func(c mqtt.Client) {
		logTmp(log).Info(ctx, tag, "Connected to MQTT Broker", domainmodel.LogMeta{"url": config.MQTTBrokerURL})
	}
	mqttOpts.OnConnectionLost = func(c mqtt.Client, err error) {
		logTmp(log).Error(ctx, tag, "Disconnected from MQTT Broker", domainmodel.LogMeta{"error": err.Error()})
	}
	mqttOpts.OnReconnecting = func(c mqtt.Client, options *mqtt.ClientOptions) {
		logTmp(log).Info(ctx, tag, "Reconnecting to MQTT Broker...", domainmodel.LogMeta{"url": config.MQTTBrokerURL})
	}
	mqttClient := mqtt.NewClient(mqttOpts)
	if tkn := mqttClient.Connect(); tkn.Wait() && tkn.Error() != nil {
		logTmp(log).Error(ctx, tag, "Failed to connect to MQTT Broker", domainmodel.LogMeta{"error": tkn.Error().Error()})
		return tkn.Error()
	}
	logTmp(log).Info(ctx, tag, "MQTT Client started successfully", nil)

	// GIN Engine
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	_ = ginEngine.SetTrustedProxies(nil)
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     config.CorsAllowedOrigins,
		AllowMethods:     config.CorsAllowedMethods,
		AllowHeaders:     config.CorsAllowedHeaders,
		AllowCredentials: config.CorsAllowCredentials,
	}))
	logTmp(log).Info(ctx, tag, "GIN Engine initated", nil)

	// Inject
	c.infrastructure = &Infrastructure{
		log:        log,
		sqlDT:      sqlDT,
		sqlTX:      sqlTX,
		sqrD:       &sqrD,
		sqrQ:       &sqrQ,
		mqttClient: mqttClient,
		ginEngine:  ginEngine,
	}
	logTmp(log).Info(ctx, tag, "Infrastructures created successfully", nil)
	return nil
}
