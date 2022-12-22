package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
)

func GetDB(ctx context.Context, connectionString string, lg *zap.Logger) (*ent.Client, *sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*15)
	defer cancelFunc()

forloop:
	for {
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("failed to wait for sql connection: %w", err)
		default:
			if err := db.Ping(); err != nil {
				lg.Warn("failed to ping db connection", zap.Error(err))
				continue forloop
			}

			break forloop
		}
	}

	entClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))

	return entClient, db, nil
}
