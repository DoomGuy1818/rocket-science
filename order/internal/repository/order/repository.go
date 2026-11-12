package order

import (
	"github.com/jackc/pgx/v5"
)

type repository struct {
	orders *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *repository {
	return &repository{
		orders: conn,
	}
}

func (r *repository) Copy() *pgx.ConnConfig {
	return r.orders.Config()
}
