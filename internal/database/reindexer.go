package database

import (
	"fmt"
	"os"

	"main/internal/models"

	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

func ConnectReindexer(REINDEXER_URL string) (*reindexer.Reindexer) {
	reindexerDB := reindexer.NewReindex(REINDEXER_URL, reindexer.WithCreateDBIfMissing())
	if err := reindexerDB.Status().Err; err != nil {
		fmt.Println("Error connecting to Reindex: ", err)
		os.Exit(-1)
	}
	// Проверка и создание коллекции persons, в случае, если она отсутствует
	err := reindexerDB.OpenNamespace(
		"persons",
		reindexer.DefaultNamespaceOptions(),
		models.Person{},
	)
	if err != nil {
		fmt.Println("Error creating namespace \"persons\": ", err)
		os.Exit(-1)
	}
	return reindexerDB
}