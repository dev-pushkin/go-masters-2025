package config

import "strconv"

type Config struct {
	API        APIConfig        `yaml:"api"`
	Logger     LoggerConfig     `yaml:"logger"`
	Metric     MetricConfig     `yaml:"metric"`
	DB         DBConfig         `yaml:"db"`
	RateReview RateReviewConfig `yaml:"rateReview"`
	Jobs       JobsConfig       `yaml:"jobs"` // Configuration for scheduled jobs
}
type APIConfig struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	ReadTimeout  int    `yaml:"readTimeout"`  // in seconds
	WriteTimeout int    `yaml:"writeTimeout"` // in seconds
	IdleTimeout  int    `yaml:"idleTimeout"`  // in seconds
}
type LoggerConfig struct {
	Level      string `yaml:"level"` // e.g., "debug", "info", "warn", "error"
	Format     string `yaml:"format"`
	GelfUdpUrl string `yaml:"gelfUdpUrl"` // URL for GELF logging, if applicable
}
type MetricConfig struct {
	Enabled   bool   `yaml:"enabled"`   // Enable or disable metrics
	Host      string `yaml:"host"`      // Host for the metrics server
	Port      int    `yaml:"port"`      // Port for the metrics server
	Endpoint  string `yaml:"endpoint"`  // Endpoint for the metrics
	Namespace string `yaml:"namespace"` // Namespace for the metrics
	Subsystem string `yaml:"subsystem"` // Subsystem for the metrics
}

type DBConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Database       string `json:"database"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	SSLMode        string `json:"sslmode"`
	PoolSize       int    `json:"pool_size"`
	MigrationsPath string `json:"migrations_path"`
}
type JobsConfig struct {
	RateReviewInterval string `yaml:"rateReviewInterval"` // Interval for rate review in cron format, e.g., "0 0 * * *" for daily
}
type RateReviewConfig struct {
	Endpoint string `yaml:"endpoint"` // Optional system prompt for the rate review provider
}

func (d *DBConfig) PgURL() string {
	// postgres://user:pass@localhost:5432/dbname?sslmode=disable
	return "postgres://" + d.Username + ":" + d.Password + "@" + d.Host + ":" + strconv.Itoa(d.Port) + "/" + d.Database + "?sslmode=" + d.SSLMode
}

func Read() (Config, error) {
	// This function should read the configuration from a file or environment variables.
	// For simplicity, we will return a hardcoded configuration.
	return Config{
		API: APIConfig{
			Port:         8080,
			Host:         "localhost",
			ReadTimeout:  10, // in seconds
			WriteTimeout: 10, // in seconds
			IdleTimeout:  60, // in seconds
		},
		Logger: LoggerConfig{
			Level:      "debug",
			Format:     "text",
			GelfUdpUrl: "",
		},
		Metric: MetricConfig{
			Enabled:   false,
			Host:      "localhost",
			Port:      9090,
			Endpoint:  "/metrics",
			Namespace: "review_manager",
			Subsystem: "api",
		},
		DB: DBConfig{
			Host:           "localhost",
			Port:           5432,
			Username:       "postgres",
			Password:       "root",
			Database:       "review_manager",
			SSLMode:        "disable",
			PoolSize:       10,
			MigrationsPath: "./review_manager/migrations",
		},
		RateReview: RateReviewConfig{
			Endpoint: "http://localhost:11434/api/generate",
		},
		Jobs: JobsConfig{
			RateReviewInterval: "@every 10s", // Daily at midnight
		},
	}, nil
}
