package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var BASE_URL = "https://www.saramin.co.kr/zf_user/search/recruit?searchType=search&searchword=react"

func main() {
	getPages()
}

func getPages() int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", BASE_URL, nil)
	checkErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")

	res, err := client.Do(req)
	checkErr(err)
	checkCode(res)
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)

	fmt.Println("=== .pagination 검색 결과 ===")
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		fmt.Println(html)
	})
	
	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	fmt.Println("Status Code:", res.StatusCode)
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status:", res.StatusCode)
	}
}