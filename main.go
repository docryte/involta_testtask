package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	// Reindexer
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"

	// Echo для API
	"github.com/labstack/echo/v4"
)

type Achievement struct {
	Content	string		`reindex:"content" json:"content" binding:"required"`
	Date	time.Time	`reindex:"date" json:"date"`
}

type Job struct {
	StartedAt		time.Time		`reindex:"started_at" json:"started_at" binding:"required"`
	EndedAt			time.Time		`reindex:"ended_at" json:"ended_at" binding:"gtfield=StartedAt"`
	Name 			string			`reindex:"name" json:"name" binding:"required"`
	Type 			string			`reindex:"type" json:"type" binding:"required"`
	Position 		string 			`reindex:"position" json:"position" binding:"required"`
	DismissalReason	string			`reindex:"dissmisal_reason" json:"dismissal_reason"`
	Achievements	[]Achievement	`reindex:"achievements" json:"achievements"`
}

type PersonPost struct {
	FirstName	string		`reindex:"first_name" json:"first_name" binding:"required"`
	LastName	string 		`reindex:"last_name" json:"last_name" binding:"required"`
	Username	string		`reindex:"username,hash" json:"username" binding:"required"`
	Birthdate	time.Time	`reindex:"birthdate" json:"birthdate" binding:"required"`
	Profession	string		`reindex:"profession" json:"profession" binding:"required"`
	Jobs		[]Job 		`reindex:"jobs" json:"jobs" binding:"required"`
}

type Person struct {
	ID 			int64 		`reindex:"id,,pk" json:"_id"`
	PersonPost
	CreatedAt 	time.Time	`reindex:"created_at" json:"created_at"`
	UpdatedAt	time.Time	`reindex:"updated_at" json:"updated_at"`
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
	}
	// Создание и запуск роутера для API и CRUD
	router := echo.New()
	router.GET("/persons", getPersons)
	router.POST("/persons", postPerson)
	router.GET("/persons/:id", getPerson)
	router.PATCH("/persons/:id", updatePerson)
	router.DELETE("/persons/:id", deletePerson)
	router.Logger.Fatal(router.Start(os.Getenv("API_ADDR")))
}

func getPersons(c echo.Context) error {
	p := 0; l := 10
	if err := echo.QueryParamsBinder(c).Int("page", &p).Int("limit", &l).BindError(); err != nil {
		return c.String(400, err.Error())
	}
	result, err := db.Query("persons").Limit(l).Offset(p*l).Exec().FetchAll()
	if err != nil {
		return c.String(500, fmt.Sprint("Error retrieving data: ", err))
	}
	return c.JSON(200, result)
}

func postPerson(c echo.Context) error {
	return c.String(500, "Doesn't Work")
	/*p := new(PersonPost)
	if err := c.Bind(&p); err != nil { return err }
	person := Person{
		0,
		*p,
		time.Now().UTC(),
		time.Now().UTC(),
	}
	if _, err := db.Insert("persons", person); err != nil {
		return c.String(500, fmt.Sprint("Error creating new person: ", err))
	}
	return c.JSON(200, person)*/
}

func getPerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, fmt.Sprint("Bad ID: ", id))
	}
	result, found := db.Query("persons").Where("_id", reindexer.EQ, id).Get()
	if !found {
		return c.String(404, "Not Found")
	}
	return c.JSON(200, result) 
}

func updatePerson(c echo.Context) error {
	return c.String(500, "Doesn't Work")
}

func deletePerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, fmt.Sprint("Bad ID: ", id))
	}
	if _, err := db.Query("persons").Where("_id", reindexer.EQ, id).Delete(); err != nil {
		return c.String(500, fmt.Sprint("Error deleting document by ID: ", err.Error()))
	}
	return c.String(200, "Success")
}