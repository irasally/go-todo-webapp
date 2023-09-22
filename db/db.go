package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var sqldb *sql.DB
var Connection *bun.DB

func CreateDBConnection(){
	var err error
	sqldb, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	Connection = bun.NewDB(sqldb, pgdialect.New())
	Connection.AddQueryHook(bundebug.NewQueryHook(
		//bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
}

func CloseDBConnection(){
	Connection.Close()
	sqldb.Close()
}
