package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-health" {
		port := env("SERVER_PORT", "3000")
		resp, err := http.Get("http://127.0.0.1:" + port + "/")
		if err != nil || resp.StatusCode != 200 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	table := os.Getenv("MYSQL_TABLE")
	port := env("SERVER_PORT", "3000")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbName)
	query := fmt.Sprintf("SELECT 1 FROM `%s` LIMIT 1", table)

	for {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println(err)
			time.Sleep(750 * time.Millisecond)
			continue
		}

		if err = db.Ping(); err != nil {
			if me, ok := err.(*mysql.MySQLError); ok {
				switch me.Number {
				case 1049:
					log.Printf("Unknown Database: %s", dbName)
				default:
					log.Println(me.Message)
				}
			} else {
				msg := err.Error()
				switch {
				case strings.Contains(msg, "Access denied"):
					log.Printf("Access denied to database '%s' with user %s", dbName, user)
				case strings.Contains(msg, "connection refused"):
					log.Printf("Server not running yet: %s", host)
				case strings.Contains(msg, "no such host"):
					log.Printf("No container with name '%s' is discoverable", host)
				default:
					log.Println(msg)
				}
			}
			db.Close()
			time.Sleep(750 * time.Millisecond)
			continue
		}

		log.Println("Connection established, database selected")

		if _, err = db.Exec(query); err != nil {
			if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1146 {
				log.Printf("No table with name %s", table)
			} else {
				log.Println(err)
			}
			db.Close()
			time.Sleep(750 * time.Millisecond)
			continue
		}

		db.Close()
		log.Println("table found")
		break
	}

	log.Printf("server listening on port %s", port)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "NOOP")
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
