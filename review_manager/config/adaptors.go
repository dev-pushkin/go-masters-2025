package config

import "github.com/go_course_master/review_manager/pkg/pgx"

func (c *Config) PgxConfig() pgx.Config {
	return pgx.Config{
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		Database: c.DB.Database,
		User:     c.DB.Username,
		Password: c.DB.Password,
		SSLMode:  c.DB.SSLMode,
		PoolSize: c.DB.PoolSize,
	}
}
