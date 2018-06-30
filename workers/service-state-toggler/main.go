package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	checkErr(err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			s := string(d.Body[:])
			user := gjson.GetMany(s, "user.id", "user.wid", "user.action")
			log.Printf("parsed: %s", s)

			resp, err := resty.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Authorization", Env["WISPRO_TOKEN"]).
				SetBody(`{"state":"` + user[2].String() + `"}`).
				Put(Env["WISPRO_BASEURL"] + user[1].String())
			checkErr(err)
			log.Printf("response: %s", resp)
			stmt, err := db.Prepare("update Users set serviceState=? where idUsers=?")
			checkErr(err)
			res, err := stmt.Exec(user[2].String(), user[0].String())
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println(affect)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
