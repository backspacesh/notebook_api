package sqlstore_test

import (
	"os"
	"testing"
)

var(
	databaseString string
)

func TestMain(m *testing.M) {
	databaseString = os.Getenv("DATABASE_URL")
	if databaseString == "" {
		databaseString = "host=localhost dbname=notebook_api user=postgres password=qwerty sslmode=disable"
	}

	os.Exit(m.Run())
}
