package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"gopkg.in/square/go-jose.v2/json"
)

func setupRabbitMQ(b []byte) {
	log.Println("I am in Line STEP 3")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	if err != nil {
		log.Printf("Error connecting RabbitMQ client: %s", err)
		log.Panic(err)
	}

	log.Printf("Successfully connected to our RabbitMQ Instance\n")

	ch, err := conn.Channel()

	if err != nil {
		log.Printf("RabbitMQ channel connection error: %s", err)
		panic(err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("RabbitMQ QueueDeclare error: %s", err)
		panic(err)
	}
	fmt.Println(q)

	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		},
	)
	log.Println("I am in Line STEP 4")
	if err != nil {
		log.Printf("RabbitMQ Publishing message error: %s", err)
		panic(err)
	}

	log.Println("I am in Line STEP 5")
}

type Book struct {
	Name      string `bson:"name" form:"name" binding:"required,min=3"`
	Author    string `bson:"author" form:"author" binding:"required,min=3"`
	PageCount int    `bson:"page_count" form:"count" binding:"required,min=1"`
}

type Author struct {
	FullName string `bson:"full_name"`
}

// PublishBook adds an album from JSON received in the request body.
func PublishBook(c *gin.Context) {
	var b Book
	if err := c.Bind(&b); err != nil {
		log.Println(b.Name)
		log.Println(b.Author)
		log.Println(b.PageCount)
		log.Printf("%#v", err)
		//c.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
	}

	log.Println("I am in Line STEP 1")

	if err := c.BindJSON(&b); err != nil {
		log.Println("I am here error : ", err)
	}
	//var bookByte []Book
	//bookByte = append(bookByte, b)
	bookData, _ := json.Marshal(b)
	setupRabbitMQ(bookData)
	//c.IndentedJSON(http.StatusCreated, "Successfully Published Message to Queue")
	log.Printf("I am here : %s", string(bookData))
	log.Println("I am in Line STEP 2")
}
func main() {
	fmt.Println("Go RabbitMQ")
	router := gin.Default()
	router.POST("/albums", PublishBook)
	router.Run("localhost:8200")
}
