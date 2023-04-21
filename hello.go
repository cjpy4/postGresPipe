package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "10.31.4.166"
	port     = 5432
	user     = "postgres"
	password = "pgpassword"
	dbname   = "goapi"
)

type Row struct {
	Id         int64  `json:"id"`
	First_name string `json:"firstName"`
	Last_name  string `json:"lastName"`
	Email      string `json:"email"`
}

func getRows() ([]Row, error) {
	var values []Row

	rows, err := db.Query("SELECT * FROM mock_data")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var value Row
		if err != nil {
			log.Println(err)
		}

		if err := rows.Scan(&value.Id, &value.First_name, &value.Last_name, &value.Email); err != nil {
			log.Fatal(err)
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return values, nil

}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		panic(pingErr)
	}

	fmt.Println("Successfully connected!")

	rows, err := getRows()

	if err != nil {
		log.Fatal(err)
	}

	obj, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(obj))

	posturl := "https://webhook.site/565e147a-604c-45dc-92de-17735c261772"

	body := obj

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

}
