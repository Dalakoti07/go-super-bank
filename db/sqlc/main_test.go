package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simplebank/utils"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

// main entry point for all tests
func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("config %v and %v \n", config.DBDriver, config.DBSource)
	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to DB: ", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
