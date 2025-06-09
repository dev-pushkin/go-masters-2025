package pgx

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib" // Importing stdlib for compatibility with pgxpool
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
	PoolSize int    `json:"pool_size"`
}

func (c *Config) connectionString() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" dbname=" + c.Database +
		" user=" + c.User +
		" password=" + c.Password +
		" sslmode=" + c.SSLMode
}

func New(config Config) (*pgxpool.Pool, error) {
	connString := config.connectionString()
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(config.PoolSize)

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
