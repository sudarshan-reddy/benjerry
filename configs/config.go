package configs

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

//Config holds all the configs needed for `qs-gateway`
type Config struct {
	LogLevel   string `envconfig:"LOG_LEVEL" required:"true"`
	LogFormat  string `envconfig:"LOG_FORMAT" required:"true"`
	Hostname   string `envconfig:"HOSTNAME" required:"true"`
	ListenPort string `envconfig:"LISTEN_PORT" required:"true"`

	PostgresDBURL            string `envconfig:"POSTGRES_DB_URL" required:"true"`
	PostgresDBMaxConnections int    `envconfig:"POSTGRES_DB_MAX_CONNECTIONS" default:"6"`
	MigrationsPath           string `envconfig:"DB_MIGRATIONS_PATH" required:"true"`
	LoadData                 bool   `envconfig:"LOAD_FIRST_TIME_DATA" required:"true"`

	StaticTokens StaticTokens `envconfig:"STATIC_TOKENS" required:"true"`
}

//Load loads all the configs
func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("BENJERRY", &config)
	return &config, err
}

//StaticTokens is a custom config type
type StaticTokens map[string][]string

//Decode implements Decoder to be able to be unmarshalled correctly
func (st *StaticTokens) Decode(value string) error {
	staticTokens := map[string][]string{}
	for _, staticToken := range strings.Split(value, ";") {
		tokenAndScopes := strings.Split(staticToken, "=")
		if len(tokenAndScopes) != 2 {
			return fmt.Errorf("invalid static token : %s", staticToken)
		}
		if _, ok := staticTokens[tokenAndScopes[0]]; ok {
			return fmt.Errorf("duplicate bearer token : %s", tokenAndScopes[0])
		}
		staticTokens[tokenAndScopes[0]] = strings.Split(tokenAndScopes[1], ",")
	}
	*st = staticTokens
	return nil
}
