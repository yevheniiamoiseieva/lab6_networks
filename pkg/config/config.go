package config

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Configuration struct {
	Application ApplicationConfiguration
	Database    DatabaseConfiguration
}

type ApplicationConfiguration struct {
	Port            string
	RequestTimeout  int
	ResponseTimeout int
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfiguration() *Configuration {
	_ = godotenv.Load()
	v := viper.New()
	v.AutomaticEnv()

	cfg := new(Configuration)

	cfg.Application.Port = v.GetString("APP_PORT")
	cfg.Application.RequestTimeout = v.GetInt("REQUEST_TIMEOUT")
	cfg.Application.ResponseTimeout = v.GetInt("RESPONSE_TIMEOUT")

	cfg.Database.Host = v.GetString("DB_HOST")
	cfg.Database.Port = v.GetInt("DB_PORT")
	cfg.Database.User = v.GetString("DB_USER")
	cfg.Database.Password = v.GetString("DB_PASSWORD")
	cfg.Database.Name = v.GetString("DB_NAME")
	cfg.Database.SSLMode = v.GetString("DB_SSLMODE")

	return cfg
}

func NewPostgresDB(cfg DatabaseConfiguration) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return db, nil
}

func RunMigrations(cfg DatabaseConfiguration) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
