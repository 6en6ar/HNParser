package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)


var hook_url = "ENTER-SLACK-CHANNEL-HOOK"

func sendToSlack(url string) {
	if !CheckLink(url) {
		return
	}
	text := map[string]string{"text": url}

	jsonValue, _ := json.Marshal(text)

	resp, err := http.Post(hook_url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	WriteToFile(url)

}
func WriteToFile(link string) {
	f, err := os.OpenFile("linkcheck.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()

	_, err2 := f.WriteString(link + "\n")
	if err2 != nil {
		log.Fatalln(err2)
	}

}
func CheckLink(link string) bool {
	f, err := os.OpenFile("linkcheck.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("[ - ] Can't open file --> " + err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == link {
			return false
		}
	}
	return true

}
func main() {
	res, err := http.Get("https://thehackernews.com/")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a.story-link").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		fmt.Println("[ + ] Sending ---> " + link)
		sendToSlack(link)
		time.Sleep(15 * time.Second)
		fmt.Println()
	})
}
