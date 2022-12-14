package movie

import (
	"context"
	"errors"

	"github.com/gcapodicimeli/api-movies/internal/domain"
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Movie, error)
	GetAllMoviesByGenre(ctx context.Context, genreID int) ([]domain.Movie, error)
	GetMovieWithContext(ctx context.Context, id int) (movie domain.Movie, err error)
	GetMovieByID(ctx context.Context, id int) (domain.Movie, error)
	Save(ctx context.Context, b domain.Movie) (domain.Movie, error)
	Update(ctx context.Context, b domain.Movie, id int) (domain.Movie, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Movie, error) {
	movies, err := s.repo.GetAll(ctx)
	if err != nil {
		return []domain.Movie{}, err
	}
	return movies, err
}

func (s *service) GetAllMoviesByGenre(ctx context.Context, genreID int) ([]domain.Movie, error) {
	movies, err := s.repo.GetAllMoviesByGenre(ctx, genreID)
	if err != nil {
		return []domain.Movie{}, err
	}
	return movies, err
}

func (s *service) GetMovieByID(ctx context.Context, id int) (movie domain.Movie, err error) {
	movie, err = s.repo.GetMovieByID(ctx, id)
	if err != nil {
		return domain.Movie{}, err
	}
	return movie, nil
}

func (s *service) GetMovieWithContext(ctx context.Context, id int) (movie domain.Movie, err error) {
	movie, err = s.repo.GetMovieWithContext(ctx, id)
	if err != nil {
		return domain.Movie{}, err
	}
	return movie, nil

}

func (s *service) Save(ctx context.Context, m domain.Movie) (domain.Movie, error) {
	if s.repo.Exists(ctx, m.ID) {
		return domain.Movie{}, errors.New("error: movie id already exists")
	}
	movie_id, err := s.repo.Save(ctx, m)
	if err != nil {
		return domain.Movie{}, err
	}

	m.ID = int(movie_id)
	return m, nil
}

func (s *service) Update(ctx context.Context, b domain.Movie, id int) (domain.Movie, error) {

	err := s.repo.Update(ctx, b, id)
	if err != nil {
		return domain.Movie{}, err
	}
	updated, err := s.repo.GetMovieByID(ctx, id)
	if err != nil {
		return b, err
	}
	return updated, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
