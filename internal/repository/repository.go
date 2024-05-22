package repository

import "movies-backend/internal/models"

type Repository interface {
	AllMovies() ([]*models.Movie, error)
	GetMovie(id int) (*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
}
