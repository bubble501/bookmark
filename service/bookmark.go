package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bubble501/bookmark/models"
	"github.com/goware/urlx"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/html/charset"
)

//GetBookmark generated an bookmark based on the given url.
func GetBookmark(bookmark *models.Bookmark) error {
	response, err := http.Get(bookmark.URL)
	if err != nil {
		fmt.Printf("%#v, %T", err, err)
		return err
	}
	defer response.Body.Close()

	utfBody, err := charset.NewReader(response.Body, "")
	if err != nil {
		fmt.Printf("%#v, %T", err, err)
		return err
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Printf("%#v, %T", err, err)
		return err
	}

	bookmark.Thumbnail, err = FetchFavicon(bookmark.URL, doc)
	if err != nil {
		fmt.Printf("%#v, %T", err, err)
		return err
	}

	bookmark.Title = getTitle(doc)
	fmt.Printf("The title is: %s\n", bookmark.Title)
	return nil
}

// FetchFavicon fetch favicon based on the description in
// http://stackoverflow.com/questions/5119041/how-can-i-get-a-web-sites-favicon
func FetchFavicon(url string, doc *goquery.Document) (string, error) {
	domain, _ := getDomain(url)
	faviconPath, _ := getFaviconPath(doc)
	if strings.HasPrefix(faviconPath, "http") == false {
		faviconPath = domain + faviconPath
	}

	fmt.Println(faviconPath)
	response, err := http.Get(faviconPath)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	filename := fmt.Sprintf("%s.ico", uuid.NewV4())
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Println("Success")
	return filename, nil
}

//Fetch form <link>
//<link rel="shortcut icon" href="/favicon.ico" />
//<link rel="icon" href="/favicon.png" />
func getFaviconPath(doc *goquery.Document) (iconPath string, err error) {
	first := doc.Find(`link[rel*="icon"]`).First()
	if first != nil {
		iconPath, _ = first.Attr("href")
	} else {
		iconPath = "/favicon.ico"
	}
	return
}

func getTitle(doc *goquery.Document) (title string) {
	first := doc.Find("title").First()
	if first != nil {
		title = first.Text()
	}
	return
}

func getDomain(url string) (domain string, ok bool) {
	u, err := urlx.Parse(url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("The host is %s \n", u.Host)
	if u.Scheme != "" {
		ok = true
		domain = u.Scheme + "://" + u.Host
	} else {
		ok = true
		domain = u.Host
	}
	fmt.Printf("The domain is %s\n", domain)
	return
}
