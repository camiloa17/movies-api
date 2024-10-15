package dbrepo

import (
	"context"
	"log"
	"movies-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *storageRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id, first_name, last_name,
			email, password, created_at,
			updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := s.PGDB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (s *storageRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id, first_name, last_name,
			email, password, created_at,
			updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := s.PGDB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
