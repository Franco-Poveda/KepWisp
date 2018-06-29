package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type jbody struct {
	State string `json:"state,omitempty"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@/KepWisp01?charset=utf8")
	checkErr(err)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"toggler", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			s := string(d.Body[:])

			user := gjson.GetMany(s, "user.id", "user.wid", "user.action")
			log.Printf("parsed: %s", s)

			resp, err := resty.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Authorization", "Token token=acba205b14cf69a6e14b461d26c0ffce43d09dc782254dc9c091ef97b413710093da8eea28e880ab21b09c41ee949f03").
				SetBody(`{"state":"` + user[2].String() + `"}`).
				//SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
				Put("http://demo8579073.mockable.io/api/contracts/" + user[1].String())
			if err != nil {
				log.Printf("parsed: %s", err)
			}
			log.Printf("parsed: %s", resp)
			// update
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
