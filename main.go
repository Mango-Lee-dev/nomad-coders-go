package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id        string
	title     string
	date      string
	condition string
	sector    string
}

var BASE_URL = "https://www.saramin.co.kr/zf_user/search/recruit?searchType=search&searchword=react"

// 전역 HTTP Client (connection pooling 활용)
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	var allJobs []extractedJob
	c := make(chan []extractedJob, 10) // buffered channel

	totalPages := getPages()
	if totalPages == 0 {
		fmt.Println("No pages found")
		return
	}

	fmt.Printf("Found %d pages to scrape\n", totalPages)

	// Rate limiting: 동시 요청 수 제한 (최대 5개)
	semaphore := make(chan struct{}, 5)

	// 모든 goroutine을 먼저 시작
	for i := 1; i <= totalPages; i++ {
		go func(page int) {
			semaphore <- struct{}{}        // 슬롯 획득
			defer func() { <-semaphore }() // 슬롯 반환
			getPage(page, c)
		}(i)
	}

	// 결과 수신
	for i := 1; i <= totalPages; i++ {
		extractedJobs := <-c
		allJobs = append(allJobs, extractedJobs...)
	}

	fmt.Println("Total jobs found:", len(allJobs))
	writeJobs(allJobs)
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	defer file.Close() // 파일 닫기 추가

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Date", "Condition", "Sector"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.date, job.condition, job.sector}
		w.Write(jobSlice)
	}

	fmt.Println("Jobs written to jobs.csv")
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".job_tit>a").Text())
	date := cleanString(card.Find(".job_date").Text())
	condition := cleanString(card.Find(".job_condition").Text())
	sector := cleanString(card.Find(".job_sector").Text())

	return extractedJob{
		id:        id,
		title:     title,
		date:      date,
		condition: condition,
		sector:    sector,
	}
}


func cleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func getPage(page int, mainC chan<- []extractedJob) {
	var extractedJobs []extractedJob
	pageUrl := BASE_URL + "&recruitPage=" + strconv.Itoa(page) + "&recruitSort=relation&recruitPageCount=40&inner_com_type=&company_cd=0%2C1%2C2%2C3%2C4%2C5%2C6%2C7%2C9%2C10&show_applied=&quick_apply=&except_read=&ai_head_hunting=&mainSearch=n"

	res, err := makeRequest(pageUrl)
	if err != nil {
		fmt.Printf("Error fetching page %d: %v\n", page, err)
		mainC <- extractedJobs // 빈 슬라이스 반환
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("Error parsing page %d: %v\n", page, err)
		mainC <- extractedJobs
		return
	}

	searchCards := doc.Find(".item_recruit")
	searchCards.Each(func(i int, s *goquery.Selection) {
		job := extractJob(s) // 동기 처리
		extractedJobs = append(extractedJobs, job)
	})

	fmt.Printf("Page %d: found %d jobs\n", page, len(extractedJobs))
	mainC <- extractedJobs
}

func getPages() int {
	pages := 0

	res, err := makeRequest(BASE_URL)
	if err != nil {
		log.Fatalln("Failed to get pages:", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// makeRequest: HTTP 요청 공통 함수
func makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, fmt.Errorf("request failed with status: %d", res.StatusCode)
	}

	return res, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
