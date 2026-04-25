package config

import (
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/utils"
)

var (
	AppHost        string                  = "0.0.0.0"
	AppPort        int                     = 80
	AppEnvironment domainmodel.Environment = domainmodel.EnvironmentDevelopment
	AppLogLevel    domainmodel.LogLevel    = domainmodel.LogLevelInfo

	CorsAllowedOrigins   []string = []string{"http://localhost:3030"}
	CorsAllowedMethods   []string = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	CorsAllowedHeaders   []string = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	CorsAllowCredentials bool     = true

	TokenIssuer            string        = "gelora-api"
	TokenAccessPrivateKey  string        = "changeme_private_key"
	TokenAccessTTL         time.Duration = 30 * time.Minute
	TokenRefreshPrivateKey string        = "refresh_private_key"
	TokenRefreshTTL        time.Duration = 72 * time.Hour

	MQTTBrokerURL      string = "tcp://localhost:1883"
	MQTTBrokerClientID string = "gelora-api"
	MQTTBrokerUsername string = "gelora-mqtt"
	MQTTBrokerPassword string = "gelora-mqtt"

	SQLiteDatabase string = "./gelora.db"
)

func LoadEnv() error {
	err := utils.EnvLoad(".env")
	if err != nil {
		return err
	}

	AppHost, err = utils.EnvGetString("APP_HOST", &AppHost)
	if err != nil {
		return err
	}

	AppPort, err = utils.EnvGetInt("APP_PORT", &AppPort)
	if err != nil {
		return err
	}

	AppEnvironment, err = utils.EnvGetEnvironment("APP_ENVIRONMENT", &AppEnvironment)
	if err != nil {
		return err
	}

	AppLogLevel, err = utils.EnvGetLogLevel("APP_LOG_LEVEL", &AppLogLevel)
	if err != nil {
		return err
	}

	CorsAllowedOrigins, err = utils.EnvGetStrings("CORS_ALLOWED_ORIGINS", CorsAllowedOrigins)
	if err != nil {
		return err
	}

	CorsAllowedMethods, err = utils.EnvGetStrings("CORS_ALLOWED_METHODS", CorsAllowedMethods)
	if err != nil {
		return err
	}

	CorsAllowedHeaders, err = utils.EnvGetStrings("CORS_ALLOWED_HEADERS", CorsAllowedHeaders)
	if err != nil {
		return err
	}

	CorsAllowCredentials, err = utils.EnvGetBool("CORS_ALLOW_CREDENTIALS", &CorsAllowCredentials)
	if err != nil {
		return err
	}

	TokenIssuer, err = utils.EnvGetString("TOKEN_ISSUER", &TokenIssuer)
	if err != nil {
		return err
	}

	TokenAccessPrivateKey, err = utils.EnvGetString("TOKEN_ACCESS_PRIVATE_KEY", &TokenAccessPrivateKey)
	if err != nil {
		return err
	}

	TokenAccessTTL, err = utils.EnvGetDuration("TOKEN_ACCESS_TTL_MINUTES", &TokenAccessTTL, time.Minute)
	if err != nil {
		return err
	}

	TokenRefreshPrivateKey, err = utils.EnvGetString("TOKEN_REFRESH_PRIVATE_KEY", &TokenRefreshPrivateKey)
	if err != nil {
		return err
	}

	TokenRefreshTTL, err = utils.EnvGetDuration("TOKEN_REFRESH_TTL_HOURS", &TokenRefreshTTL, time.Hour)
	if err != nil {
		return err
	}

	MQTTBrokerURL, err = utils.EnvGetString("MQTT_BROKER_URL", &MQTTBrokerURL)
	if err != nil {
		return err
	}

	MQTTBrokerClientID, err = utils.EnvGetString("MQTT_BROKER_CLIENT_ID", &MQTTBrokerClientID)
	if err != nil {
		return err
	}

	MQTTBrokerUsername, err = utils.EnvGetString("MQTT_BROKER_USERNAME", &MQTTBrokerUsername)
	if err != nil {
		return err
	}

	MQTTBrokerPassword, err = utils.EnvGetString("MQTT_BROKER_PASSWORD", &MQTTBrokerPassword)
	if err != nil {
		return err
	}

	SQLiteDatabase, err = utils.EnvGetString("SQLITE_DATABASE", &SQLiteDatabase)
	if err != nil {
		return err
	}

	return nil
}
