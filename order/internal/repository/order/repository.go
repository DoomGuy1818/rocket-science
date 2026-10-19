package order

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	orders *pgx.Conn
}

func NewStorage(ctx context.Context, dbURI string) *repository {
	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		os.Exit(1)
	}

	return &repository{
		orders: conn,
	}
}

func (r *repository) Copy() *pgx.ConnConfig {
	return r.orders.Config()
}
