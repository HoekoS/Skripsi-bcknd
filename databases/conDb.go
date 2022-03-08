package databases

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect() (dbPool *pgxpool.Pool, err error) {
	// "postgres://postgres:salatigahatiberiman@34.101.155.53:5432/postgres"
	// dbPool, err = pgxpool.Connect(context.Background(), os.Getenv("GO_DATABASE_URL"))
	dbPool, err = pgxpool.Connect(context.Background(), "postgres://postgres:salatiga123@34.101.251.192/postgres")

	if err != nil {
		return nil, err
	}

	return
}
