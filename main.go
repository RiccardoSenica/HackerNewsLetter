package main

import (
	"context"
	"flag"
	"fmt"
	"hackernewsletter/db"
	"hackernewsletter/hackernews"
	"hackernewsletter/params"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	conf, err := params.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	table := flag.String("t", "", "The name of the table")
	flag.Parse()

	if *table == "" {
		fmt.Println("Table name not specified. Using default name.")
		table = &conf.TableName
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(conf.AwsRegion))
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.DescribeTableInput{
		TableName: table,
	}

	newsTable := db.Table{client, *table, conf.TableIndex}

	resp, err := db.GetTableInfo(context.TODO(), client, input)
	if err != nil {
		println(("Table not found. Creating it..."))
		_, new_err := db.CreateTable(newsTable)
		if new_err != nil {
			panic("Failed creating table " + *table)
		}

		resp, err = db.GetTableInfo(context.TODO(), client, input)
	}

	fmt.Printf("Table %v has %v elements\n", *table, resp.Table.ItemCount)

	newsIds := hackernews.GetTopNewsIds(conf.TopNews)

	var newsBatch []db.News

	for i := 0; i < 10; i++ {
		myNews := hackernews.GetNewsById(newsIds[i], conf.SingleNews)
		newsBatch = append(newsBatch, myNews)
	}

	db.AddNewsBatch(newsTable, newsBatch, conf.BatchSize)

	timeEnd := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	timeStart := timeEnd.Add(-time.Hour * 24)

	todayNews, _ := db.ReadTodayNews(newsTable, int(timeStart.Unix()), int(timeEnd.Unix()))
	println("NEWS", todayNews)
}
