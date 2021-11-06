package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streadway/amqp"
)

func GetResponse() {
	fmt.Println("Consumber Application")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	defer conn.Close()
	if err != nil {
		log.Printf("Error connecting RabbitMQ client: %s", err)
		log.Panic(err)
	}

	ch, err := conn.Channel()
	defer ch.Close()

	if err != nil {
		log.Printf("Error RabbitMQ Channel : %s", err)
		log.Panic(err)
	}

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			msg := string(d.Body)
			fmt.Printf("Recieved Message: %s\n", msg)
			var b Book
			json.Unmarshal(d.Body, &b)
			addPost(b)
			//json.Unmarshal(msg, &b)
			//log.Printf("%#v, %T", d.Body, d.Body)
		}
	}()

	log.Printf("Successfully connected to our RabbitMQ Instance\n")
	fmt.Println(" [*] - waiting for messages ")

	<-forever
}

type Book struct {
	Name      string `bson:"name" form:"name" binding:"required,min=3"`
	Author    string `bson:"author" form:"author" binding:"required,min=3"`
	PageCount int    `bson:"page_count" form:"count" binding:"required,min=1"`
}

type Author struct {
	FullName string `bson:"full_name"`
}

func addPost(b Book) {

	//_name := fmt.Sprintf("Vijay Sandesh %s", msg)
	//params := url.Values{"name": string(b.Name), "author": {"Ramdhari sing Dinkar"}, "count": {"13"}}
	params := url.Values{"name": {b.Name}, "author": {b.Author}, "count": {strconv.Itoa(b.PageCount)}}
	resp, err := http.PostForm("http://localhost:8100/albums",
		params)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	// Log the request body
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		print(err)
	}
	fmt.Println(data["form"])
	readMessase()

}
func readMessase() {
	resp, err := http.Get("http://localhost:8100/albums")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
func main() {
	GetResponse()
}
