package main

import (
	"fmt"
	"os"
	"time"

	// Reindexer
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"

	// Для валидации json в post/put
	"github.com/go-playground/validator"

	// Echo для API
	"github.com/labstack/echo/v4"
)

type Achievement struct {
	Content	string		`reindex:"content" json:"content" validate:"required"`
	Date	time.Time	`reindex:"date" json:"date"`
}

type Job struct {
	StartedAt		time.Time		`reindex:"started_at" json:"started_at" validate:"required"`
	EndedAt			time.Time		`reindex:"ended_at" json:"ended_at" validate:"gtfield=StartedAt"`
	Name 			string			`reindex:"name" json:"name" validate:"required"`
	Type 			string			`reindex:"type" json:"type" validate:"required"`
	Position 		string 			`reindex:"position" json:"position" validate:"required"`
	DismissalReason	string			`reindex:"dissmisal_reason" json:"dismissal_reason"`
	Achievements	[]Achievement	`reindex:"achievements" json:"achievements"`
}

type PersonPost struct {
	FirstName	string		`reindex:"first_name" json:"first_name" validate:"required"`
	LastName	string 		`reindex:"last_name" json:"last_name" validate:"required"`
	Username	string		`reindex:"username,hash" json:"username" validate:"required"`
	Birthdate	time.Time	`reindex:"birthdate" json:"birthdate" validate:"required"`
	Profession	string		`reindex:"profession" json:"profession" validate:"required"`
	Jobs		[]Job 		`reindex:"jobs" json:"jobs" validate:"required"`
}

type Person struct {
	ID 			int64 		`reindex:"id,,pk" json:"_id"`
	PersonPost
	UpdatedAt	time.Time	`reindex:"updated_at" json:"updated_at"`
}

// Структуры для валидации запросов с json
type PlaygroundValidator struct {
	validator *validator.Validate
}

func (pv *PlaygroundValidator) Validate(i interface{}) error {
	if err := pv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}

var db *reindexer.Reindexer

func main() {
	// Подключение к Reindexer по ссылке из окружения и проверка подключения
	db = reindexer.NewReindex(os.Getenv("RX_CONNECTION_URL"), reindexer.WithCreateDBIfMissing())
	if err := db.Status().Err; err != nil {
		fmt.Println("Error connecting to Reindex: ", err)
		os.Exit(-1)
	}
	// Проверка и создание коллекции persons, в случае, если она отсутствует
	err := db.OpenNamespace("persons", reindexer.DefaultNamespaceOptions(), Person{})
	if err != nil {
		fmt.Println("Error creating namespace \"persons\": ", err)
		os.Exit(-1)
	}
	// Создание и запуск роутера для API и CRUD
	router := echo.New()
	router.Validator = &PlaygroundValidator{validator: validator.New()}
	router.GET("/persons", getPersons)
	router.POST("/persons", postPerson)
	router.GET("/persons/:id", getPerson)
	router.PUT("/persons/:id", putPerson)
	router.DELETE("/persons/:id", deletePerson)
	router.Logger.Fatal(router.Start(os.Getenv("API_ADDR")))
}

func getPersons(c echo.Context) error {
	p := 0; l := 10
	if err := echo.QueryParamsBinder(c).Int("page", &p).Int("limit", &l).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	result, err := db.Query("persons").Limit(l).Offset(p*l).Exec().FetchAll()
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprint("Error retrieving Person: ", err))
	}
	return c.JSON(200, result)
}

func postPerson(c echo.Context) error {
	p := new(PersonPost)
	if err := c.Bind(p); err != nil { 
		return echo.NewHTTPError(400, err.Error()) 
	}
	if err := c.Validate(p); err != nil {
		return err
	}
	person := Person{
		0,
		*p,
		time.Now().UTC(),
	}
	_, err := db.Insert(
		"persons",
		&person,
		"ID=serial()",
	)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprint("Error creating Person: ", err.Error()))
	}
	return c.JSON(200, person)
}

func getPerson(c echo.Context) error {
	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	result, found := db.Query("persons").Where("_id", reindexer.EQ, id).Get()
	if !found {
		return c.String(404, "Not Found")
	}
	return c.JSON(200, result) 
}

func putPerson(c echo.Context) error {
	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	p := new(PersonPost)
	if err := c.Bind(p); err != nil { 
		return echo.NewHTTPError(400, err.Error()) 
	}
	person := Person{
		id,
		*p,
		time.Now().UTC(),
	}
	_, err := db.Update(
		"persons",
		&person,
	)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprint("Error updating Person: ", err.Error()))
	}
	return c.JSON(200, person)
}

func deletePerson(c echo.Context) error {
	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := db.Delete("persons", &Person{ID: int64(id)}); err != nil {
		return c.String(500, fmt.Sprint("Error deleting document by ID: ", err.Error()))
	}
	return c.String(200, "Success")
}