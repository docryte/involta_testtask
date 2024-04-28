package main

import (
	"fmt"
	"os"

	// Reindexer
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type Achievement struct {
	Content	string	`reindex:"content"`
	Date	int64	`reindex:"date"`
}

type Job struct {
	StartedAt		int64			`reindex:"started_at"`
	EndedAt			int64			`reindex:"ended_at"`
	Name 			string			`reindex:"name"`
	Type 			string			`reindex:"type"`
	Position 		string 			`reindex:"position"`
	DismissalReason	string			`reindex:"dissmisal_reason"`
	Achievements	[]Achievement	`reindex:"achievements"`
}

type Person struct {
	ID 			int64 		`reindex:"id,,pk"`
	FirstName	string		`reindex:"first_name"`
	SecondName	string 		`reindex:"second_name"`
	Username	string		`reindex:"username,,pk"`
	Birthdate	int64		`reindex:"birthdate"`
	Profession	string		`reindex:"profession"`
	Jobs		[]Job 		`reindex:"jobs"`
	CreatedAt 	int64		`reindex:"created_at"`
	UpdatedAt	int64		`reindex:"updated_at"`
}

func main() {
	// Подключение к Reindexer по ссылке из окружения и проверка подключения
	db := reindexer.NewReindex(os.Getenv("RX_CONNECTION_URL"), reindexer.WithCreateDBIfMissing())
	if err := db.Status().Err; err != nil {
		fmt.Println("Error connecting to Reindex: ", err)
		os.Exit(-1)
	}
	// Проверка и создание коллекции persons, в случае, если она отсутствует
	err := db.OpenNamespace("persons", reindexer.DefaultNamespaceOptions(), Person{})
	if err != nil {
		fmt.Println("Error creating namespace \"persons\": ", err)
	}
}