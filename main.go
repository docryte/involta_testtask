package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	// Reindexer
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"

	// Для валидации json в post/put
	"github.com/go-playground/validator"

	// Echo для API
	"github.com/labstack/echo/v4"

	// Redis для кэша
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/cache/v9"
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
	Sort 			int				`reindex:"sort" json:"sort" validate:"required"`
}

type PersonPost struct {
	FirstName	string		`reindex:"first_name" json:"first_name" validate:"required"`
	LastName	string 		`reindex:"last_name" json:"last_name" validate:"required"`
	Username	string		`reindex:"username,hash" json:"username" validate:"required"`
	Birthdate	time.Time	`reindex:"birthdate" json:"birthdate" validate:"required"`
	Profession	string		`reindex:"profession" json:"profession" validate:"required"`
	Jobs		[]Job 		`reindex:"jobs" json:"jobs" validate:"required"`
}

/*
Структура документа с двойной вложенностью, как в задании

Разделение на две структуры, потому что я не приудмал, 
как лучше гарантировать, что пользователь в теле json не передаст _id или updated_at

И в доках echo есть что-то похожее:
https://echo.labstack.com/docs/binding#security
*/
type Person struct {
	ID 			int64 		`reindex:"id,,pk" json:"_id"`
	PersonPost
	UpdatedAt	time.Time	`reindex:"updated_at" json:"updated_at"`
}

/* 
Структуры для валидации запросов с json

В доках echo:
https://echo.labstack.com/docs/request#validate-data
*/
type PlaygroundValidator struct {
	validator *validator.Validate
}

func (pv *PlaygroundValidator) Validate(i interface{}) error {
	if err := pv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}

var reindexer_db *reindexer.Reindexer
var redis_cache *cache.Cache

func main() {
	// Подключение к Reindexer по ссылке из окружения и проверка подключения
	reindexer_db = reindexer.NewReindex(os.Getenv("REINDEXER_URL"), reindexer.WithCreateDBIfMissing())
	if err := reindexer_db.Status().Err; err != nil {
		fmt.Println("Error connecting to Reindex: ", err)
		os.Exit(-1)
	}
	// Проверка и создание коллекции persons, в случае, если она отсутствует
	err := reindexer_db.OpenNamespace(
		"persons", 
		reindexer.DefaultNamespaceOptions(), 
		Person{},
	)
	if err != nil {
		fmt.Println("Error creating namespace \"persons\": ", err)
		os.Exit(-1)
	}

	// Подключение к Redis и создание клиента для кэша
	database, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		fmt.Printf("REDIS_DATABASE environment variable should be int, not: %v", os.Getenv("REDIS_DATABASE"))
		os.Exit(-1)
	}
	redis_db := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       database,
    })
	redis_cache = cache.New(&cache.Options{
		Redis: redis_db,
	})

	// Создание и запуск роутера для API и CRUD
	router := echo.New()
	router.Validator = &PlaygroundValidator{validator: validator.New()}
	router.GET("/persons", getPersons)
	router.POST("/persons", postPerson)
	router.GET("/persons/:id", getPerson)
	router.PUT("/persons/:id", putPerson)
	router.DELETE("/persons/:id", deletePerson)
	router.Logger.Fatal(router.Start(os.Getenv("API_URL")))
}

func getPersons(c echo.Context) error {
	p := 0; l := 10
	if err := echo.QueryParamsBinder(c).Int("page", &p).Int("limit", &l).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	result, err := reindexer_db.Query("persons").Limit(l).Offset(p*l).Exec().FetchAll()
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
	_, err := reindexer_db.Insert(
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
	person := new(Person)
	// Сначала пробуем получить документ из кэша
	err := redis_cache.Get(context.Background(), string(rune(id)), person)
	/* 
	Если документа нет в кэше, метод вернёт ошибку ErrCacheMiss
	Ошибку обрабатываем для получения документа из БД и его записи в кэш
	
	Возможно, имеет смысл сначала вызывать Exists для проверки кэша на наличие
	*/
	if err == cache.ErrCacheMiss {
		result, found := reindexer_db.Query("persons").Where("_id", reindexer.EQ, id).Get()
		if !found {
			return echo.NewHTTPError(404, "Not Found")
		}
		person = result.(*Person)
		redis_cache.Set(&cache.Item{
			Ctx: 	context.Background(),
			Key: 	string(rune(id)),
			Value: 	person,
			// Объект в кэше хранится ровно 15 минут	
			TTL:	time.Duration(time.Duration.Minutes(15)),
		})
	}
	return c.JSON(200, person) 
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
	_, err := reindexer_db.Update(
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
	if err := reindexer_db.Delete("persons", &Person{ID: int64(id)}); err != nil {
		return c.String(500, fmt.Sprint("Error deleting document by ID: ", err.Error()))
	}
	return c.String(200, "Success")
}