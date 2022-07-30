package handler

import (
	"bytes"
	"context"
	"fmt"
	"hackernewsletter/db"
	"hackernewsletter/hackernews"
	"hackernewsletter/mail"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Handler(ctx context.Context) (string, error) {
	batchSize, _ := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	fetchSize, _ := strconv.Atoi(os.Getenv("FETCH_SIZE"))

	var newsBatch []db.News
	table := hackernews.HackernewsTable()
	htmlTemplate := hackernews.HackernewsTemplate()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("Unable to load SDK config, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.DescribeTableInput{
		TableName: &table,
	}

	newsTable := db.Table{client, table}

	resp, err := db.GetTableInfo(context.TODO(), client, input)
	if err != nil {
		println(("Table not found. Creating it..."))
		_, new_err := db.CreateTable(newsTable)
		if new_err != nil {
			log.Fatalln("Failed creating table.")
		}

		resp, err = db.GetTableInfo(context.TODO(), client, input)
	}

	fmt.Printf("Table %v has %v elements\n", table, resp.Table.ItemCount)

	newsIds := hackernews.GetTopNewsIds("https://hacker-news.firebaseio.com/v0/topstories.json")

	for i := 0; i < fetchSize; i++ {
		myNews := hackernews.GetNewsById(newsIds[i], "https://hacker-news.firebaseio.com/v0/item/{ID}.json")
		newsBatch = append(newsBatch, myNews)
	}

	db.AddNewsBatch(newsTable, newsBatch, batchSize)

	timeEnd := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	timeStart := timeEnd.Add(-time.Hour * 24)

	todayNews, _ := db.ReadTodayNews(newsTable, int(timeStart.Unix()), int(timeEnd.Unix()))

	t, err := template.ParseFiles(fmt.Sprintf("assets/%v.gohtml", htmlTemplate))
	if err != nil {
		log.Fatalln(err)
	}

	emailNews := struct {
		News []db.News
	}{
		News: todayNews,
	}

	var emailNewsBuffer bytes.Buffer

	err = t.Execute(&emailNewsBuffer, emailNews)
	if err != nil {
		log.Fatalln(err)
	}

	mail.SendNewsletter(emailNewsBuffer.String())

	return "Completed.", err
}
