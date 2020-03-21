package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createData = `
	CREATE TABLE data(
		id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, 
		created DATETIME DEFAULT CURRENT_TIMESTAMP,
		model  TEXT
	)
	`
	insertData = `
	INSERT INTO
		data (model) 
	VALUES 
		($1)
	`

	selectData = `
	SELECT 
		created, 
		model
	FROM
		data
	WHERE
		id = $1
	`
)

func execWithContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) error {
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func queryRowWithContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) error {
	var created sql.NullTime
	var model sql.NullString

	if err := db.QueryRowContext(ctx, query, args...).Scan(&created, &model); err != nil {
		return err
	}
	if created.Valid {
		fmt.Println(created)
	}
	if model.Valid {
		fmt.Println(model)
	}
	return nil

}

func main() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	defer os.Remove("data.db")
	defer db.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := execWithContext(ctx, db, createData); err != nil {
		panic(err)
	}
	log.Println("table 'Data' created")

	for i := 0; i < 100; i++ {
		if err := execWithContext(ctx, db, insertData, fmt.Sprintf("model-%d", i)); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond)
	}
	log.Println("100 rows inserted")
	if err := queryRowWithContext(ctx, db, selectData, 5); err != nil {
		panic(err)
	}

}
