package movie

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gcapodicimeli/api-movies/internal/domain"
	"github.com/stretchr/testify/assert"
)

var null_int *int

var movie_test = domain.Movie{
	ID:           1,
	Created_at:   time.Now(),
	Updated_at:   time.Now(),
	Title:        "Cars 1",
	Rating:       4,
	Awards:       2,
	Release_date: time.Layout,
	Length:       0,
	Genre_id:     0,
}

func TestGetOneWithContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(movie_test.ID, movie_test.Title, movie_test.Rating, movie_test.Awards, movie_test.Length, movie_test.Genre_id)
	mock.ExpectQuery(regexp.QuoteMeta(GET_MOVIE)).WithArgs(movie_test.ID).WillReturnRows(rows)

	repo := NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movieResult, err := repo.GetMovieWithContext(ctx, movie_test.ID)
	assert.NoError(t, err)
	assert.Equal(t, movie_test.Title, movieResult.Title)
	assert.Equal(t, movie_test.ID, movieResult.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExist_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(movie_test.ID)
	mock.ExpectQuery(regexp.QuoteMeta(EXIST_MOVIE)).WithArgs(movie_test.ID).WillReturnRows(rows)

	repo := NewRepository(db)
	result := repo.Exists(context.TODO(), movie_test.ID)

	assert.True(t, result)
}

func TestGetAll_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(column)
	movies := []domain.Movie{{ID: 1, Title: "Avatar", Rating: 22, Awards: 99, Length: 0, Genre_id: 1}, {ID: 2, Title: "Simpson", Rating: 33, Awards: 11, Length: 2, Genre_id: 2}}

	for _, m := range movies {
		rows.AddRow(m.ID, m.Title, m.Rating, m.Awards, m.Length, m.Genre_id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_MOVIES)).WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.GetAll(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, movies, result)
}
