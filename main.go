package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var BASE_URL = "https://kr.indeed.com/jobs?q=python&limit=50&vjk=211d68e47e6d8074"

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(BASE_URL)
	checkErr(err)
	checkCode(res)
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)

	doc.Find(".navigation").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
	
	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {

}