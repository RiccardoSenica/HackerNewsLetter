package main

import (
	"context"
	"flag"
	"fmt"
	"hackernewsletter/db"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Config struct {
	TABLE_NAME string
}

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

	new_table := db.Table{client, *table}

	resp, err := db.GetTableInfo(context.TODO(), client, input)
	if err != nil {
		println(("Table not found. Creating it..."))
		_, new_err := db.CreateTable(new_table)
		if new_err != nil {
			panic("Failed creating table " + *table)
		}

		resp, err = db.GetTableInfo(context.TODO(), client, input)
	}

	fmt.Printf("Table %v has %v elements", *table, resp.Table.ItemCount)

	var news []db.News

	news = append(news, db.News{Id: 1, Title: "First", Url: "www.first.com"})
	news = append(news, db.News{Id: 2, Title: "Second", Text: "Second text"})

	db.AddNewsBatch(new_table, news)
}
