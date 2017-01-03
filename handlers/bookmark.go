package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/bubble501/bookmark/logger"
	"github.com/bubble501/bookmark/models"
	"github.com/bubble501/bookmark/service"
	"github.com/labstack/echo"
)

var log = logger.Logger

//H is an map of response.
type H map[string]interface{}

// GetBookmarks get the bookmark.
func GetBookmarks(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Infoln("GetBookmarks has been invoked in handler package")
		return c.JSON(http.StatusOK, models.GetBookmarks(db))
	}
}

// PutBookmark put the bookmark.
func PutBookmark(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Instantiate a new task
		var bookmark models.Bookmark
		// Map imcoming JSON body to the new Task
		c.Bind(&bookmark)
		if strings.HasPrefix(bookmark.URL, "http://") == false &&
			strings.HasPrefix(bookmark.URL, "https://") == false {
			bookmark.URL = "http://" + bookmark.URL
		}
		log.Infoln("PutBookmark has been invoked in handler package")
		duplicateURL, err := models.IsURLExist(db, bookmark.URL)
		if err != nil {
			return nil
		}
		res := map[string]int{"go": 1}
		if duplicateURL == true {
			log.Warningln("Duplicate url ", bookmark.URL)
			return c.JSON(http.StatusOK, H{
				"error":  "网址已存在！",
				"result": res,
			})
		}
		// Add a task using our new model

		job := service.FetchAndSaveJob{Bookmark: &bookmark, Db: db}
		service.TaskQueue.AddJob(job)
		id, err := models.PutBookmark(db, bookmark.URL)

		// Return a JSON response if successful
		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
				"hello":   res,
			})
		}
		log.Errorln("models.PutBookmark failed with error ", err)
		return err

	}
}

// DeleteBookmark  endpoint
func DeleteBookmark(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		// Use our new model to delete a task
		_, err := models.DeleteBookmark(db, id)
		// Return a JSON response on success
		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})
		}
		log.Errorln("models.DeleteBookmark failed with error ", err)
		return err

	}
}
