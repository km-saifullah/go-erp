package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/km-saifullah/go-erp/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// Maximum number of open connections.
	db.SetMaxOpenConns(25)

	// Maximum number of idle connections.
	db.SetMaxIdleConns(10)

	// Maximum amount of time a connection can be reused.
	db.SetConnMaxLifetime(30 * time.Minute)

	// Maximum amount of time an idle connection can remain idle.
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Verify that MySQL is actually reachable.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, nil
}
