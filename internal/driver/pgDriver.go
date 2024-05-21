package driver

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type PgStorage struct {
	DB *pgx.Conn
}

func NewPgStorage(dsn string) (*PgStorage, error) {

	dbCon, err := connectToDb(dsn)
	if err != nil {
		return &PgStorage{}, err
	}
	return &PgStorage{
		DB: dbCon,
	}, nil
}

func connectionError(e error) {
	log.Println(e)
}

func connectToDb(dsn string) (*pgx.Conn, error) {
	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		connectionError(err)
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		connectionError(err)
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		connectionError(err)
		return nil, err
	}
	return conn, nil
}
