package pkg

import (
	"database/sql"
	"fmt"

	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     utils.GetEnvOrDefault("DATABASE_HOST", "localhost"),
		Port:     utils.GetEnvOrDefault("DATABASE_PORT", "5432"),
		Username: utils.GetEnvOrDefault("DATABASE_USERNAME", "postgres"),
		Password: utils.GetEnvOrDefault("DATABASE_PASS", ""),
		Name:     utils.GetEnvOrDefault("DATABASE_NAME", "ptm_db"),
		SSLMode:  utils.GetEnvOrDefault("DATABASE_SSLMODE", "disable"),
	}
}

func (c *Config) buildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func (c *Config) DSN() string {
	return c.buildDSN()
}

func Connect() error {
	return ConnectWithConfig(NewConfig())
}
func ConnectWithConfig(config *Config) error {
	dsn := config.buildDSN()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
