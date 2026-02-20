// Package env has helper functions to deal with environment variables.
// This package uses viper for env management.
package env

import (
	"errors"
	"fmt"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ErrMissingDSN is returned when neither MYSQL_DSN nor DB_DSN is set.
var ErrMissingDSN = errors.New("config: MYSQL_DSN or DB_DSN is required")

// Environment variable key constants (viper keys; map to uppercase env vars).
const (
	Port        = "port"
	HTTPPort    = "http_port"
	MySQLDSN    = "mysql_dsn"
	DBDSN       = "db_dsn"
	HealthToken = "xdrv_healthz_token"
	HealthzToken = "healthz_token"
	Version     = "version"
	BuildInfo   = "build_info"
	StartTime   = "start_time"
)

// MandatoryParameters lists env vars that must be set.
// DSN (MYSQL_DSN or DB_DSN) is validated in config.Load.
var MandatoryParameters = []string{}

// OptionalParameters holds optional env vars with their defaults.
var OptionalParameters = map[string]interface{}{
	Port:        ":8080",
	HTTPPort:    "",
	HealthToken: "",
	HealthzToken: "",
	Version:     "",
	BuildInfo:   "",
	StartTime:   "",
}

// InitEnv populates all env parameters in viper.
// envPrefix is set for automatic env binding (e.g. service name).
// If mandatory params are not set, returns an error.
// Populates default values for optional params.
// After init, params can be accessed via viper.Get or viper.GetString.
func InitEnv(envPrefix string, mandatoryParams []string, optionalParams map[string]interface{}) (skipped []string, err error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)

	var notPresent []string
	for _, v := range mandatoryParams {
		if !viper.IsSet(v) {
			notPresent = append(notPresent, v)
		}
	}
	if len(notPresent) != 0 {
		msg := fmt.Sprintf("mandatory env variables not set: %+v", notPresent)
		return nil, pkgerrors.New(msg)
	}

	for k, v := range optionalParams {
		if v == nil {
			skipped = append(skipped, k)
			continue
		}
		viper.SetDefault(k, v)
	}
	return skipped, nil
}

// GetDSN returns MySQL DSN (tries MYSQL_DSN then DB_DSN). Returns ErrMissingDSN if neither set.
func GetDSN() (string, error) {
	dsn := viper.GetString(MySQLDSN)
	if dsn == "" {
		dsn = viper.GetString(DBDSN)
	}
	if dsn == "" {
		return "", ErrMissingDSN
	}
	return dsn, nil
}

// GetPort returns HTTP listen address with ":" prefix if needed.
func GetPort() string {
	port := viper.GetString(Port)
	if port == "" {
		port = viper.GetString(HTTPPort)
	}
	if port == "" {
		port = ":8080"
	}
	if len(port) > 0 && port[0] != ':' {
		port = ":" + port
	}
	return port
}

// GetHealthToken returns healthz token (tries XDRV_HEALTHZ_TOKEN then HEALTHZ_TOKEN).
func GetHealthToken() string {
	token := viper.GetString(HealthToken)
	if token == "" {
		token = viper.GetString(HealthzToken)
	}
	return token
}

// GetString returns trimmed viper string for key.
func GetString(key string) string {
	return strings.TrimSpace(viper.GetString(key))
}
