package handler

import (
	"bytes"
	"context"
	"fmt"
	"hackernewsletter/db"
	"hackernewsletter/hackernews"
	"hackernewsletter/mail"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Handler(ctx context.Context) (string, error) {
	table := "news_table"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
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
			panic("Failed creating table " + table)
		}

		resp, err = db.GetTableInfo(context.TODO(), client, input)
	}

	fmt.Printf("Table %v has %v elements\n", table, resp.Table.ItemCount)

	newsIds := hackernews.GetTopNewsIds("https://hacker-news.firebaseio.com/v0/topstories.json")

	var newsBatch []db.News

	for i := 0; i < 10; i++ {
		myNews := hackernews.GetNewsById(newsIds[i], "https://hacker-news.firebaseio.com/v0/item/{ID}.json")
		newsBatch = append(newsBatch, myNews)
	}

	db.AddNewsBatch(newsTable, newsBatch, 25)

	timeEnd := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	timeStart := timeEnd.Add(-time.Hour * 24)

	todayNews, _ := db.ReadTodayNews(newsTable, int(timeStart.Unix()), int(timeEnd.Unix()))

	t, err := template.ParseFiles("mail/index.gohtml")

	if err != nil {
		panic(err)
	}

	data := struct {
		News []db.News
	}{
		News: todayNews,
	}

	var doc bytes.Buffer

	err = t.Execute(&doc, data)

	if err != nil {
		panic(err)
	}

	body := doc.String()

	mail.SendNewsletter(body)

	return "Completed.", err
}
