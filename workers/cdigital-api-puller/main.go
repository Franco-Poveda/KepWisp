package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

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

	c, err := NewConsumer(Env["RABBIT_URI"], "topic_logs", "topic", "q0", "info", "api-pooler")
	if err != nil {
		log.Fatalf("%s", err)
	}
	lifetime := flag.Duration("lifetime", 0, "lifetime of process before shutdown (0s=infinite)")
	if *lifetime > 0 {
		log.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {

	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	var Env map[string]string
	Env, err := godotenv.Read()
	db, err := sql.Open("mysql", Env["MYSQL_URI"])
	checkErr(err)

	for d := range deliveries {

		log.Printf("Received a message: %s", d.Body)

		//t, err := time.Parse("20060102", "20180704")
		//for now := time.Now(); t.Before(now); t = t.AddDate(0, 0, 1) {
		t := time.Now()
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
		d.Ack(true)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}
