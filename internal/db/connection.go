package db

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/config"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
)

func GetDB(cfg config.DB) (*ent.Client, *sql.DB, error) {
	db, err := sql.Open("pgx", cfg.GetConnectionString())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping sql connection: %w", err)
	}

	opts := []ent.Option{
		ent.Driver(entsql.OpenDB(dialect.Postgres, db)),
		// ent.Debug(),
	}

	if cfg.ShowSql {
		opts = append(opts, ent.Debug())
	}

	entClient := ent.NewClient(opts...)

	return entClient, db, nil
}

// SELECT o.*, os.* FROM orders o JOIN LATERAL (SELECT * FROM order_status WHERE order_order_status = o.id AND DATE(current_date) = (SELECT MAX(DATE(current_date)) FROM order_status WHERE order_order_status = o.id) ORDER BY id DESC LIMIT 1) os ON TRUE WHERE os.order_status_name_order_status = 'In Review';
