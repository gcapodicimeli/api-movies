package movie

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gcapodicimeli/api-movies/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	null_int     *int
	ERRORFORZADO = errors.New("Error forzado")
)

var movie_test = domain.Movie{
	ID:           1,
	Created_at:   time.Now(),
	Updated_at:   time.Now(),
	Title:        "Cars 1",
	Rating:       4,
	Awards:       2,
	Release_date: time.Layout,
	Length:       null_int,
	Genre_id:     null_int,
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
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(column)
	movies := []domain.Movie{{ID: 1, Title: "Avatar", Rating: 22, Awards: 99, Length: nil, Genre_id: nil}, {ID: 2, Title: "Simpson", Rating: 33, Awards: 11, Length: nil, Genre_id: nil}}

	for _, m := range movies {
		rows.AddRow(m.ID, m.Title, m.Rating, m.Awards, m.Length, m.Genre_id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_MOVIES)).WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.GetAll(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, movies, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(column)
	movies := []domain.Movie{{ID: 1, Title: "Avatar", Rating: 22, Awards: 99, Length: nil, Genre_id: nil}, {ID: 2, Title: "Simpson", Rating: 33, Awards: 11, Length: nil, Genre_id: nil}}

	for _, m := range movies {
		rows.AddRow(m.ID, m.Title, m.Rating, m.Awards, m.Length, m.Genre_id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_MOVIES)).WillReturnError(ERRORFORZADO) // * Acá es la diferencia para que falle

	repo := NewRepository(db)

	result, err := repo.GetAll(context.TODO())

	assert.EqualError(t, err, ERRORFORZADO.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(column)
	movie := domain.Movie{ID: 1, Title: "Avatar", Rating: 22, Awards: 99, Length: nil, Genre_id: nil}

	rows.AddRow(movie.ID, movie.Title, movie.Rating, movie.Awards, movie.Length, movie.Genre_id)

	mock.ExpectQuery(regexp.QuoteMeta(GET_MOVIE)).WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.GetMovieByID(context.TODO(), movie_test.ID)

	assert.NoError(t, err)
	assert.Equal(t, movie, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_MOVIE))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_MOVIE)).WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(movie_test.ID, movie_test.Title, movie_test.Rating, movie_test.Awards, movie_test.Length, movie_test.Genre_id)
	// mock.ExpectQuery(regexp.QuoteMeta(GET_MOVIE)).WithArgs(1).WillReturnRows(rows)

	repository := NewRepository(db)
	ctx := context.TODO()

	newID, err := repository.Save(ctx, movie_test)
	assert.NoError(t, err)

	// movieResult, err := repository.GetMovieByID(ctx, int(newID))
	assert.NoError(t, err)

	// assert.NotNil(t, movieResult)
	// assert.Equal(t, movie_test.ID, movieResult.ID)
	assert.Equal(t, int64(movie_test.ID), newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_MOVIE))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_MOVIE)).WithArgs(movie_test.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	result := repository.Delete(context.TODO(), int64(movie_test.ID))

	assert.NoError(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	movie := domain.Movie{ID: 1, Title: "Avatar", Rating: 22, Awards: 99, Length: nil, Genre_id: nil}

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_MOVIE))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_MOVIE)).WithArgs(movie.Title, movie.Rating, movie.Awards, movie.Length, movie.Genre_id, movie.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{"id", "title", "rating", "awards", "length", "genre_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(movie.ID, movie.Title, movie.Rating, movie.Awards, movie.Length, movie.Genre_id)

	repository := NewRepository(db)
	result := repository.Update(context.TODO(), movie, movie.ID)

	assert.NoError(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
