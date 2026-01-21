package feature_postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/k0kubun/pp"
)

type BookModel struct {
	ID              int
	Title           string
	Author          string
	Review          *string
	PublicationYear time.Time
	IsRead          bool
	AddedAt         time.Time
	FinishedAt      *time.Time
}

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, "postgres://postgres:1234@localhost:5432/hwdb")
}

func CreateTableBooks(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(200) NOT NULL,
	author VARCHAR(100) NOT NULL,
	review VARCHAR(200),
	publication_year TIMESTAMP NOT NULL,
	is_read BOOLEAN NOT NULL,
	added_at TIMESTAMP NOT NULL,
	finished_at TIMESTAMP
	);
	`

	if _, err := conn.Exec(ctx, sqlQuery); err != nil {
		return err
	}
	return nil
}

func InsertBook(ctx context.Context, conn *pgx.Conn, book BookModel) error {
	sqlQuery := `
	INSERT INTO books (title, author, review, publication_year, is_read, added_at, finished_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := conn.Exec(ctx, sqlQuery, book.Title, book.Author, book.Review, book.PublicationYear, book.IsRead, book.AddedAt, book.FinishedAt); err != nil {
		return err
	}

	return nil
}

func SelectBooks(ctx context.Context, conn *pgx.Conn) ([]BookModel, error) {
	sqlQuery := `
	SELECT * FROM books
	ORDER BY id ASC
	`
	books := make([]BookModel, 0)
	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var book BookModel
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Review, &book.PublicationYear, &book.IsRead, &book.AddedAt, &book.FinishedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func UpdateBook(ctx context.Context, conn *pgx.Conn, book BookModel) error {
	sqlQuery := `
	UPDATE books
	SET title = $1, author = $2, review = $3, publication_year = $4, is_read = $5, added_at = $6, finished_at = $7
	WHERE id = $8
	`

	_, err := conn.Exec(ctx, sqlQuery, book.Title, book.Author, book.Review, book.PublicationYear, book.IsRead, book.AddedAt, book.FinishedAt, book.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBook(ctx context.Context, conn *pgx.Conn, booksID []int) error {
	sqlQuery := `
	DELETE FROM books
	WHERE id = ANY($1)
	`

	if _, err := conn.Exec(ctx, sqlQuery, booksID); err != nil {
		return err
	}

	return nil
}

func ListPages(ctx context.Context, conn *pgx.Conn, N int) error {
	BooksPages, err := SelectBooks(ctx, conn)
	if err != nil {
		return err
	}
	pages := (len(BooksPages) / N) + 1
	offset := 0
	for i := 1; i <= pages; i++ {
		sqlQuery := `
		SELECT * FROM books
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
		`

		books := make([]BookModel, 0)
		rows, err := conn.Query(ctx, sqlQuery, N, offset)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var book BookModel
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Review, &book.PublicationYear, &book.IsRead, &book.AddedAt, &book.FinishedAt)
			if err != nil {
				return err
			}
			books = append(books, book)
		}

		fmt.Print("Страница " + strconv.Itoa(i) + ": ")
		pp.Print(books)
		offset += N
	}
	return nil
}
