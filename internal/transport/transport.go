package transport

import (
	"main/internal/models"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/restream/reindexer"
)

var reindexerDB *reindexer.Reindexer

func StartRouting(rdb *reindexer.Reindexer) {
	reindexerDB = rdb

	router := echo.New()
	router.Validator = &models.PlaygroundValidator{Validator: validator.New()}
	router.GET("/persons", getPersons)
	router.POST("/persons", postPerson)
	router.GET("/persons/:id", getPerson)
	router.PUT("/persons/:id", putPerson)
	router.DELETE("/persons/:id", deletePerson)
	router.Logger.Fatal(router.Start(":8000"))

}