package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	"github.com/joho/godotenv"
)

func EnvLoad(paths ...string) error {
	err := godotenv.Load(paths...)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("error loading .env file(s): %w", err)
		}
	}
	return nil
}

func EnvGetString(key string, fallback *string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		if fallback != nil {
			return *fallback, nil
		}
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}

func EnvGetStrings(key string, fallback []string) ([]string, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return fallback, nil
		}
		return nil, err
	}
	return strings.Split(valStr, ","), nil
}

func EnvGetEnvironment(key string, fallback *domainmodel.Environment) (domainmodel.Environment, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return *fallback, nil
		}
		return "", err
	}
	if !slices.Contains([]string{
		string(domainmodel.EnvironmentDevelopment),
		string(domainmodel.EnvironmentStaging),
		string(domainmodel.EnvironmentProduction),
	}, valStr) {
		return "", fmt.Errorf("environment variable %s value is not valid", key)
	}
	return domainmodel.Environment(valStr), nil
}

func EnvGetLogLevel(key string, fallback *domainmodel.LogLevel) (domainmodel.LogLevel, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return *fallback, nil
		}
		return "", err
	}
	if !slices.Contains([]string{
		string(domainmodel.LogLevelError),
		string(domainmodel.LogLevelWarn),
		string(domainmodel.LogLevelInfo),
		string(domainmodel.LogLevelDebug),
	}, valStr) {
		return "", fmt.Errorf("environment variable %s value is not valid", key)
	}
	return domainmodel.LogLevel(valStr), nil
}

func EnvGetInt(key string, fallback *int) (int, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return *fallback, nil
		}
		return 0, err
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an integer", key)
	}
	return val, nil
}

func EnvGetBool(key string, fallback *bool) (bool, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return *fallback, nil
		}
		return false, err
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a boolean", key)
	}
	return val, nil
}

func EnvGetDuration(key string, fallback *time.Duration, unit time.Duration) (time.Duration, error) {
	valStr, err := EnvGetString(key, nil)
	if err != nil {
		if fallback != nil {
			return *fallback, nil
		}
		return 0, err
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an integer", key)
	}
	return time.Duration(val) * unit, nil
}
