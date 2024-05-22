package dbrepo

import (
	"context"
	"log"
	"movies-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *storageRepo) AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id, title, release_date,
			runtime, mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		FROM movies
	`
	var movies []*models.Movie
	rows, err := s.PGDB.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return movies, err
	}

	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return movies, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil

}

func (s *storageRepo) GetMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			id, title, release_date,
			runtime, mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		FROM movies
		WHERE id = $1
	`
	var movie models.Movie
	err := s.PGDB.QueryRow(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &movie, nil
}
