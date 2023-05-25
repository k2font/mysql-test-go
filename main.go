package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

type user struct {
	id        int
	username  string
	password  string
	createdAt time.Time
}

var (
	MYSQL_ID       string
	MYSQL_PASSWORD string
	MYSQL_PORT     string
	MYSQL_DBNAME   string
)

func loadEnv() (string, string, string, string) {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}

	return os.Getenv("MYSQL_ID"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DBNAME")
}

func main() {
	MYSQL_ID, MYSQL_PASSWORD, MYSQL_PORT, MYSQL_DBNAME = loadEnv()
	db, err := sql.Open("mysql", MYSQL_ID+":"+MYSQL_PASSWORD+"@(127.0.1:"+MYSQL_PORT+")/"+MYSQL_DBNAME+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// テーブルの作成
	{
		query := `
			CREATE TABLE workspaces3 (
				id INT AUTO_INCREMENT,
				username TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at DATETIME,
				PRIMARY KEY (id)
		);`

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}

	// 行の挿入
	{
		username := "k2font"
		password := "test"
		createdAt := time.Now()

		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		id, err := result.LastInsertId()
		fmt.Println(id)
	}

	// ユーザデータのクエリ
	{
		rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var users []user
		for rows.Next() {
			var u user

			err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
			if err != nil {
				log.Fatal()
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v\n", users)
	}
}
