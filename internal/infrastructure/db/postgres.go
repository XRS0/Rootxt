package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/rootix/portfolio/internal/infrastructure/db/migrations"
)

func Open(ctx context.Context) (*bun.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		pgHost := getenv("PGHOST", "localhost")
		pgPort := getenv("PGPORT", "5432")
		pgUser := getenv("PGUSER", "postgres")
		pgPassword := getenv("PGPASSWORD", "postgres")
		pgDatabase := getenv("PGDATABASE", "portfolio")
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase)
	}

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := sqldb.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := RunMigrations(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		return err
	}
	if _, err := migrator.Migrate(ctx); err != nil {
		return err
	}
	return nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
