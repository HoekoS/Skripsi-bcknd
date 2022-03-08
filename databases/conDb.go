package databases

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect() (dbPool *pgxpool.Pool, err error) {
	dbPool, err = pgxpool.Connect(context.Background(), "postgres://postgres:salatigahatiberiman@34.101.155.53:5432/postgres")

	if err != nil {
		return nil, err
	}

	return
}
