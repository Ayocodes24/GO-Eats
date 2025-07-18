package tests

import (
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"log"
	"log/slog"
	"os"
)

func Setup() database.Database {
	slog.Info("Initializing Setup..")
	testDb := database.NewTestDB()

	if err := testDb.Migrate(); err != nil {
		log.Fatalf("Error migrating database: %s", err)
	}
	return testDb
}

func Teardown(testDB database.Database) {
	err := testDB.Close()
	if err != nil {
		log.Fatalf("Error closing testDB: %s", err)
	}
	err = os.RemoveAll("./tmp")
	if err != nil {
		log.Fatalf("Error removing ./tmp: %s", err)
	}
}
