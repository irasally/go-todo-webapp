package model

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var sqldb *sql.DB
var db *bun.DB

func CreateDBConnection(){
	var err error
	sqldb, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db = bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		//bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	CreateTable(db)
}

func CloseDBConnection(){
	db.Close()
	sqldb.Close()
}

// Task型のテーブルを作成する
func CreateTable(db *bun.DB) {
	var err error
	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*Todo)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
