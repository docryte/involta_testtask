package transport

import (
	"fmt"
	"main/internal/database"
	"main/internal/models"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/restream/reindexer/v3"
)

func getPersons(c echo.Context) error {
	p := 0; l := 10
	if err := echo.QueryParamsBinder(c).Int("page", &p).Int("limit", &l).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	iterator := reindexerDB.Query("persons").Limit(l).Offset(p*l).Exec()
	if iterator.Count() == 0 {
		return echo.NewHTTPError(404, "Not Found")
	}
	var persons []*models.Person
	for iterator.Next() {
		persons = append(persons, iterator.Object().(*models.Person))
	}
	// Обратная сортировка по целочисленному полю Sort 
	go sort.SliceStable(
		persons, 
		func(i, j int) bool {
			return persons[i].Sort > persons[j].Sort
		},
	)
	return c.JSON(200, persons)
}

func postPerson(c echo.Context) error {
	p := new(models.PersonPost)
	if err := c.Bind(p); err != nil { 
		return echo.NewHTTPError(400, err.Error()) 
	}
	if err := c.Validate(p); err != nil {
		return err
	}
	person := models.Person{
		ID: 		0,
		PersonPost: *p,
		UpdatedAt: 		time.Now().UTC(),
	}
	_, err := reindexerDB.Insert(
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
	person := new(models.Person)
	// Сначала пробуем получить документ из кэша
	exists := database.GetCache(string(rune(id)), person)
	/* 
	Если документа нет в кэше, метод вернёт ошибку ErrCacheMiss
	Ошибку обрабатываем для получения документа из БД и его записи в кэш
	
	Возможно, имеет смысл сначала вызывать Exists для проверки кэша на наличие
	*/
	if !exists {
		result, found := reindexerDB.Query("persons").Where("_id", reindexer.EQ, id).Get()
		if !found {
			return echo.NewHTTPError(404, "Not Found")
		}
		person = result.(*models.Person)
		// Кэшируем объект уже после сортировки
		database.SetCache(
			string(rune(id)),
			person,
		)
	}
	return c.JSON(200, person) 
}

func putPerson(c echo.Context) error {
	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	p := new(models.PersonPost)
	if err := c.Bind(p); err != nil { 
		return echo.NewHTTPError(400, err.Error()) 
	}
	person := models.Person{
		ID: 		id,
		PersonPost: *p,
		UpdatedAt: 	time.Now().UTC(),
	}
	_, err := reindexerDB.Update(
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
	if err := reindexerDB.Delete("persons", &models.Person{ID: id}); err != nil {
		return echo.NewHTTPError(500, fmt.Sprint("Error deleting document by ID: ", err.Error()))
	}
	return c.String(200, "Success")
}