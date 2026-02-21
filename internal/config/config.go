package config

import (
	"github.com/sanketmote/go-todo-service/internal/env"
)

// ErrMissingDSN re-exports env.ErrMissingDSN for backward compatibility.
var ErrMissingDSN = env.ErrMissingDSN

// Config holds application configuration loaded from environment.
type Config struct {
	Port        string // HTTP server listen address (e.g. ":8080")
	DSN         string // MySQL connection string
	HealthToken string // Token for GET /healthz
	Version     string // Service version (optional, for healthz)
	BuildInfo   string // Build info (optional, for healthz)
	StartTime   string // Service start time (optional, for healthz)
}

// Load calls env.InitEnv with envPrefix (use "" for PORT, MYSQL_DSN, etc.).
// Returns error if DSN is empty.
func Load(envPrefix string) (*Config, error) {
	_, err := env.InitEnv(envPrefix, env.MandatoryParameters, env.OptionalParameters)
	if err != nil {
		return nil, err
	}

	dsn, err := env.GetDSN()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        env.GetPort(),
		DSN:         dsn,
		HealthToken: env.GetHealthToken(),
		Version:     env.GetString(env.Version),
		BuildInfo:   env.GetString(env.BuildInfo),
		StartTime:   env.GetString(env.StartTime),
	}, nil
}
