package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hamdysherif/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can't load config file: ", err)
	}

	log.Println("Driver is: ", config.DBDriver)

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
