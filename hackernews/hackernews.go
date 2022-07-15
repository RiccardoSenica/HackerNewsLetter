package hackernews

import (
	"encoding/json"
	"fmt"
	"hackernewsletter/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetTopNewsIds(url string) (response []string) {
	fmt.Println("Retrieving news...")

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body_string := string(body)
	response = strings.Split(body_string[1:len(body_string)-1], ",")

	fmt.Println("Done.")

	return response
}

func GetNewsById(id string, url string) (response db.News) {
	fmt.Printf("Retrieving data for news with id %v...\n", id)

	news_url := strings.ReplaceAll(url, "{ID}", id)
	res, err := http.Get(news_url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &response)

	fmt.Printf("News with id %v retrieved.\n", id)

	return response
}
