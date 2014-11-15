package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schemaUp = `CREATE TABLE schedule
(id INTEGER PRIMARY KEY AUTOINCREMENT,
	hour INTEGER NOT NULL,
	days STRING NOT NULL,
	target_temp INTEGER NOT NULL);`

var schemaDown = "DROP TABLE schedule;"

type Database struct {
	connection *sqlx.DB
}

func (d *Database) Connection() *sqlx.DB {
	if d.connection != nil {
		return d.connection
	}

	dbName := fmt.Sprintf("%s-%s.db", DB_NAME, ENV)
	db, err := sqlx.Connect("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	return db
}

func (d *Database) Create() error {
	_, err := d.Connection().Exec(schemaUp)
	return err
}

func (d *Database) Destroy() error {
	_, err := d.Connection().Exec(schemaDown)
	d.Connection().Close()
	return err
}
