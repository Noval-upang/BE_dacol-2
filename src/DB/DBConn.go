package DB

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv"
)

type DBConn struct {
	db    *sql.DB
	query string
}

func isErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initDb(query string) DBConn {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@%s/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_PROTOCOL"),
			os.Getenv("DB_NAME")))
	isErr(err)

	return DBConn{db, query}
}

func (ctx DBConn) Single() *sql.Row {
	defer ctx.db.Close()
	return ctx.db.QueryRow(ctx.query)
}

func (ctx DBConn) Many() *sql.Rows {
	defer ctx.db.Close()
	rows, err := ctx.db.Query(ctx.query)

	isErr(err)

	return rows
}

func (ctx DBConn) Exec() sql.Result {
	defer ctx.db.Close()
	result, err := ctx.db.Exec(ctx.query)

	isErr(err)

	return result
}
