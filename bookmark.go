package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bubble501/bookmark/config"
	"github.com/bubble501/bookmark/handlers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	dbpath := config.Singleton.GetStringValue("dbPath", "storage.db")
	fmt.Printf("the db path is: %s", dbpath)
	db := initDB(dbpath)
	migrate(db)

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("/bookmarks", handlers.GetBookmarks(db))
	e.PUT("/bookmark", handlers.PutBookmark(db))
	e.POST("/login", handlers.Login(db))
	e.DELETE("/bookmark/:id", handlers.DeleteBookmark(db))
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)
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
