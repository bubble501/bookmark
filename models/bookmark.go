package models

import (
	"database/sql"
	"fmt"
	//load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Bookmark is a struct containing Bookmark data
type Bookmark struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Date        string `json:"date"`
}

// BookmarkCollection is a collection of bookmark.
type BookmarkCollection struct {
	Bookmarks []Bookmark `json:"items"`
}

// GetBookmarks from the DB
func GetBookmarks(db *sql.DB) BookmarkCollection {
	sql := "SELECT id, url, title, thumbnail, description, body, date FROM bookmark"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		panic(err)
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	result := BookmarkCollection{}
	for rows.Next() {
		bookmark := Bookmark{}
		err2 := rows.Scan(&bookmark.ID,
			&bookmark.URL,
			&bookmark.Title,
			&bookmark.Thumbnail,
			&bookmark.Description,
			&bookmark.Body,
			&bookmark.Date)
		// Exit if we get an error
		if err2 != nil {
			panic(err2)
		}
		result.Bookmarks = append(result.Bookmarks, bookmark)
	}
	return result
}

//IsURLExist is used to tell if url has been added already.
func IsURLExist(db *sql.DB, url string) (bool, error) {
	var id int
	err := db.QueryRow("SELECT id FROM bookmark where url=?", url).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

// PutBookmark into DB
func PutBookmark(db *sql.DB, url string) (int64, error) {
	sql := "INSERT INTO bookmark(url) VALUES(?)"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Replace the '?' in our prepared statement with 'name'
	result, err2 := stmt.Exec(url)
	// Exit if we get an error
	if err2 != nil {
		fmt.Printf("%#v, %T", err2, err2)
		return 0, err2
	}
	return result.LastInsertId()
}

func UpdateBookmark(db *sql.DB, bookmark *Bookmark) error {
	sql := "UPDATE bookmark set title=?, thumbnail=? where url=?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(bookmark.Title, bookmark.Thumbnail, bookmark.URL)
	if err2 != nil {
		fmt.Printf("%#v, %T", err2, err2)
		return err2
	}
	return nil
}

//PutBookmarkSec put whole bookmark into DB
func PutBookmarkSec(db *sql.DB, bookmark *Bookmark) (int64, error) {
	sql := "INSERT INTO bookmark(title, url, thumbnail, description, body)" +
		"VALUES(?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Replace the '?' in our prepared statement with 'name'
	result, err2 := stmt.Exec(bookmark.Title, bookmark.URL, bookmark.Thumbnail,
		bookmark.Description, bookmark.Body)
	// Exit if we get an error
	if err2 != nil {
		fmt.Printf("%#v, %T", err2, err2)
		return 0, err2
	}
	return result.LastInsertId()
}

// DeleteBookmark from DB
func DeleteBookmark(db *sql.DB, id int) (int64, error) {
	sql := "DELETE FROM bookmark WHERE id = ?"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}

	// Replace the '?' in our prepared statement with 'id'
	result, err2 := stmt.Exec(id)
	// Exit if we get an error
	if err2 != nil {
		panic(err2)
	}

	return result.RowsAffected()
}
