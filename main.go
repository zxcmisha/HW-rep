package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn_string := os.Getenv("CONN_STRING")
	conn, err := pgx.Connect(ctx, conn_string)
	if err != nil {
		panic(err)
	}

	sqlQuery := `
	CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	full_name VARCHAR(200) NOT NULL,
	phone_number VARCHAR(100)
	);
	`
	if _, err := conn.Exec(ctx, sqlQuery); err != nil {
		panic(err)
	}

	new_user := os.Getenv("NEW_USER")
	switch new_user {
	case "YES":
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Введите полное имя: ")
		scanner.Scan()
		full_name := scanner.Text()

		fmt.Print("Введите ваш номер телефона: ")
		scanner.Scan()
		phone_number := scanner.Text()

		sqlQuery := `
		INSERT INTO users (full_name, phone_number)
		VALUES ($1, $2);
		`
		if _, err := conn.Exec(ctx, sqlQuery, full_name, phone_number); err != nil {
			panic(err)
		}
		return
	case "NO":
		sqlQuery := `
		SELECT * FROM users
		ORDER BY id ASC;
		`
		rows, err := conn.Query(ctx, sqlQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var full_name string
			var phone_number string
			if err := rows.Scan(&id, &full_name, &phone_number); err != nil {
				panic(err)
			}
			fmt.Println(id, full_name, phone_number)
		}
		return
	default:
		fmt.Println("Ошибка, нужно ввести YES или NO")
	}
}
