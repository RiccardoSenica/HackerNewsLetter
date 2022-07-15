package main

import (
	"context"
	"flag"
	"fmt"
	"hackernewsletter/db"
	"hackernewsletter/hackernews"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	table := flag.String("t", "", "The name of the table")
	flag.Parse()

	if *table == "" {
		fmt.Println("You must specify a table name (-t TABLE)")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.DescribeTableInput{
		TableName: table,
	}

	newTable := db.Table{client, *table}

	resp, err := db.GetTableInfo(context.TODO(), client, input)
	if err != nil {
		println(("Table not found. Creating it..."))
		_, new_err := db.CreateTable(newTable)
		if new_err != nil {
			panic("Failed creating table " + *table)
		}

		resp, err = db.GetTableInfo(context.TODO(), client, input)
	}

	fmt.Printf("Table %v has %v elements", *table, resp.Table.ItemCount)

	newsIds := hackernews.GetTopNewsIds()

	var newsBatch []db.News

	for i := 0; i < 10; i++ {
		myNews := hackernews.GetNewsById(newsIds[i])
		newsBatch = append(newsBatch, myNews)
	}

	db.AddNewsBatch(newTable, newsBatch)
}
