package main

import (
	"fmt"
	"os"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type Location struct {
	Country		string 		`reindex:"country"`
	City		string 		`reindex:"city"`
	Street		string 		`reindex:"street"`
	Home		int			`reindex:"home"`
	
	Entrance	int			`reindex:"entrance"`
	Floor		int			`reindex:"floor"`
}

type User struct {
	ID 			int64 		`reindex:"id,,pk"`
	Username	string		`reindex:"username,,pk"`
	FirstName 	string 		`reindex:"first_name"`
	LastName 	string 		`reindex:"last_name"`
	Location 	Location	`reindex:"location"`
}

func main() {
	db := reindexer.NewReindex(os.Getenv("RX_CONNECTION_URL"), reindexer.WithCreateDBIfMissing())
	db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), User{})
	err := db.Upsert("users", &User{
		1,
		"docryte",
		"Фёдор",
		"Акулов",
		Location{
			"Россия",
			"Иваново",
			"2-я Сосневская",
			19,
			0,
			0,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	elem, found := db.Query("users").Where("ID", reindexer.EQ, 1).Get()
	if found {
		item := elem.(*User)
		fmt.Println(*item)
	}
}