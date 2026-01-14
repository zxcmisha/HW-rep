package main

import (
	"context"
	"fmt"
	"learn/feature_postgres"
)

var review string = "Arckham 2024"

func main() {
	ctx := context.Background()
	conn, err := feature_postgres.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	if err := feature_postgres.CreateTableBooks(ctx, conn); err != nil {
		panic(err)
	}
	// book := feature_postgres.BookModel{
	// 	ID:              1,
	// 	Title:           "Robin",
	// 	Author:          "Dimka",
	// 	Review:          &review,
	// 	PublicationYear: time.Date(2024, time.November, 11, 21, 00, 00, 11, time.Local),
	// 	IsRead:          false,
	// 	AddedAt:         time.Now(),
	// }
	feature_postgres.ListPages(ctx, conn, 10)
	fmt.Println("success")

}
