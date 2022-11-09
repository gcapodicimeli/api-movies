package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	service movie.Service
}

func NewMovie(service movie.Service) *Movie {
	return &Movie{
		service: service,
	}
}

func (m *Movie) GetMovieByID() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		id, err := strconv.Atoi((ctx.Param("id")))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		movie, err := m.service.GetMovieByID(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, movie)
	}
}
