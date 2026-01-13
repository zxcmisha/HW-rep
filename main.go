package main

import (
	"context"
	"fmt"
	"learn/feature_postgres"
)

func main() {
	ctx := context.Background()
	_, err := feature_postgres.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("succes")
}
