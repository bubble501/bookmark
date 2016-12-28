package main

import (
	"database/sql"
	"fmt"

	"github.com/bubble501/bookmark/config"
	"github.com/bubble501/bookmark/handlers"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbpath := config.Singleton.GetStringValue("dbPath", "storage.db")
	fmt.Printf("the db path is: %s", dbpath)
	db := initDB(dbpath)
	migrate(db)

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("/bookmarks", handlers.GetBookmarks(db))
	e.PUT("/bookmark", handlers.PutBookmark(db))
	e.DELETE("/bookmark/:id", handlers.DeleteBookmark(db))

	e.Start(":8000")
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS bookmark(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title VARCHAR NOT NULL DEFAULT "",
		url VARCHAR NOT NULL UNIQUE,
		thumbnail VARCHAR NOT NULL DEFAULT "",
		description TEXT NOT NULL DEFAULT "",
		body TEXT NOT NULL DEFAULT "",
		date DATETIME DEFAULT CURRENT_TIME
	);
	`

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}
