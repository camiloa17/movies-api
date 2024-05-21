package dbrepo

import (
	"movies-backend/internal/repository"
	"time"

	"github.com/jackc/pgx/v5"
)

type storageRepo struct {
	PGDB *pgx.Conn
}

const dbTimeout = time.Second * 3

func NewStorageRepo(pgCon *pgx.Conn) repository.Repository {
	return &storageRepo{
		PGDB: pgCon,
	}
}
