package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"encoding/csv"
	"log"
	"net/http"
	"net/url"
)

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1

	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	var Env map[string]string
	Env, err := godotenv.Read()
	db, err := sql.Open("mysql", Env["MYSQL_URI"])
	checkErr(err)

	t, err := time.Parse("20060102", "20180704")
	for now := time.Now(); t.Before(now); t = t.AddDate(0, 0, 1) {
		u, err := url.Parse(Env["URI"])
		if err != nil {
			log.Fatal(err)
		}
		q := u.Query()
		q.Set("control", Env["KEY"])
		q.Set("fecha", t.Format("20060102"))
		q.Set("hour1", "00")
		q.Set("min1", "00")
		q.Set("hour2", "23")
		q.Set("min2", "59")

		u.RawQuery = q.Encode()
		fmt.Println(u)
		data, err := readCSVFromURL(u.String())
		if err != nil {
			panic(err)
		}

		log.Printf("parsed: %s", data)
		data = data[:len(data)-1]

		for _, d := range data {
			d[2] = strings.Replace(d[2], ".", "", -1)
			s := [7]string{}
			s[0] = "1"
			s[1] = d[0][len(d[0])-4:] + "-" + d[0][2:len(d[0])-4] + "-" + d[0][0:2] + " " + d[1][0:2] + ":" + d[1][2:4] + ":" + d[1][4:6]
			s[2] = strings.Replace(d[2], ",", ".", -1)
			s[3] = d[5]
			s[4] = d[6]
			s[5] = d[7]
			s[6] = d[8]
			log.Printf("AAR: %s", s)

			stmt, err := db.Prepare("INSERT INTO `Transactions` (`type`,`tdate`,`amount`,`barcode`,`ref`,`method`,`cduid`) VALUES (?, ?,?,?,?,?,?)")
			checkErr(err)
			_, err = stmt.Exec(s[0], s[1], s[2], s[3], s[4], s[5], s[6])
			if err != nil {
				log.Print(err)
			}
		}
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
