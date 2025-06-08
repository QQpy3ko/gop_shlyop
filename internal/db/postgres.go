package db

import (
	"context"
	"errors"
	"fmt"
	"gop_shlyop/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // driver
	"github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog/log"
)

// creates a new PostgreSQL connection pool and runs migrations.
func NewPostgresDB(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	if err := runMigrations(connString); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

func runMigrations(connString string) error {
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		// Log the error but don't fail if we just can't get the version
		if !errors.Is(err, migrate.ErrNilVersion) {
			log.Warn().Err(err).Msg("could not get migration version")
		}
	} else {
		log.Info().Uint64("version", uint64(version)).Bool("dirty", dirty).Msg("DB migrations applied")
	}

	// close the source and database connections used by migrate
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Warn().Err(srcErr).Msg("migration source close error")
	}
	if dbErr != nil {
		log.Warn().Err(dbErr).Msg("migration db close error")
	}

	return nil
}
