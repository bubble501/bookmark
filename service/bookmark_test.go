package service_test

import (
	"fmt"
	"testing"

	"github.com/bubble501/bookmark/service"
)

func TestFetchFavicon(*testing.T) {
	//service.FetchFavicon("http://www.baidu.com", 1)
	//	service.FetchFaviconFromLink("https://www.ampproject.org/docs/guides/discovery", 1)
	url := "https://www.ampproject.org/docs/guides/discovery"
	bookmark := service.GetBookmark(url)
	fmt.Printf("%v", bookmark)
}

// func TestUrlPath(*testing.T) {
// 	urls := []string{
// 		"http://www.sohu.com/hello",
// 		"https://www.sina.com.cn/abc",
// 		"www.ifeng.com",
// 	}
//
// 	for index := range urls {
// 		fmt.Println(urls[index])
// 		u, err := urlx.Parse(urls[index])
// 		if err != nil {
// 			panic(err)
// 		}
// 		var domain string
// 		fmt.Printf("The host is %s \n", u.Host)
// 		if u.Scheme != "" {
// 			domain = u.Scheme + "://" + u.Host
// 		} else {
// 			domain = u.Host
// 		}
// 		fmt.Printf("The domain is %s\n", domain)
// 	}
// }
