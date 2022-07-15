package hackernews

import (
	"encoding/json"
	"hackernewsletter/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetTopNewsIds() (response []string) {
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body_string := string(body)
	response = strings.Split(body_string[1:len(body_string)-1], ",")

	return response
}

func GetNewsById(id string) (response db.News) {
	news_url := "https://hacker-news.firebaseio.com/v0/item/" + id + ".json"
	res, err := http.Get(news_url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &response)

	return response
}
