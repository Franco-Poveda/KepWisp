package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var Env map[string]string
	Env, err := godotenv.Read()

	db, err := sql.Open("mysql", Env["MYSQL_URI"])
	checkErr(err)
	defer db.Close()

	conn, err := amqp.Dial(Env["RABBIT_URI"])
	checkErr(err)
	defer conn.Close()
	ch, err := conn.Channel()
	checkErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"toggler", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	checkErr(err)

	rows, err := db.Query(overdueClients)
	checkErr(err)

	for rows.Next() {
		var idUser int

		err = rows.Scan(&idUser)
		checkErr(err)
		fmt.Println(idUser)

		body := `{"user":{"id":` + strconv.Itoa(idUser) + `, "wid":"12","action":"disable"}}`
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		checkErr(err)
	}
}

const (
	overdueClients = "SELECT idUser FROM clientValance WHERE payedServiceUntil > NOW()"
)
