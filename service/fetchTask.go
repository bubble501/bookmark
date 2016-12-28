package service

import (
	"database/sql"
	"fmt"

	"github.com/bubble501/bookmark/models"
	"github.com/bubble501/taskQueue"
)

//TaskQueue is used to fetch and save bookmark's favicon and title.
var TaskQueue = taskQueue.New(3, 2000)

func init() {
	fmt.Print("init taskQueue")
	TaskQueue.Start()
}

//FetchAndSaveJob is a job object which is used to fetch favicon and title
// and save the bookmark to db.
type FetchAndSaveJob struct {
	Bookmark *models.Bookmark
	Db       *sql.DB
}

//Execute do the fetch and save.
func (job FetchAndSaveJob) Execute() {
	err := GetBookmark(job.Bookmark)
	if err == nil {
		models.UpdateBookmark(job.Db, job.Bookmark)
	}
}
