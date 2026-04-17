package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	date string
	condition string
	sector string
}

var BASE_URL = "https://www.saramin.co.kr/zf_user/search/recruit?searchType=search&searchword=react"

func main() {
	var allJobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		allJobs = append(allJobs, extractedJobs...)
	}

	fmt.Println("Total jobs found:", len(allJobs))
	writeJobs(allJobs)
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Date", "Condition", "Sector"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.date, job.condition, job.sector}
		w.Write(jobSlice)
	}

	fmt.Println("Jobs written to CSV")
}

func extractJobs(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-gnb_idx")
	title := card.Find(".job_tit>a").Text()
	date := card.Find(".job_date").Text()
	condition := card.Find(".job_condition").Text()
	sector := card.Find(".job_sector").Text()
	
	c <- extractedJob{id: id, title: title, date: date, condition: condition, sector: sector}
}

func getPage(page int, mainC chan<- []extractedJob) {
	var extractedJobs []extractedJob
	c := make(chan extractedJob)
	pageUrl := BASE_URL + "&recruitPage=" + strconv.Itoa(page) + "&recruitSort=relation&recruitPageCount=40&inner_com_type=&company_cd=0%2C1%2C2%2C3%2C4%2C5%2C6%2C7%2C9%2C10&show_applied=&quick_apply=&except_read=&ai_head_hunting=&mainSearch=n"

	client := &http.Client{}
	req, err := http.NewRequest("GET", pageUrl, nil)
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

	searchCards := doc.Find(".item_recruit")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJobs(s, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		extractedJobs = append(extractedJobs, job)
	}
	mainC <- extractedJobs
}

func getPages() int {
	pages := 0
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

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
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