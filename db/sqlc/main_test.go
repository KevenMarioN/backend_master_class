package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDB)
	os.Exit(m.Run())
}
