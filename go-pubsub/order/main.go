package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
)

// Order represents an order from the frontend
type Order struct {
	//how do I make bind json fail if mising field?
	//we can use required tag
	Email   string `json:"email" binding:"required"`
	Product string `json:"product" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
}

type OrderCreated struct {
	Order Order  `json:"order"`
	Id    string `json:"id"`
}

func generateOrderID() string {
	id := rand.Intn(99999)
	return fmt.Sprintf("%05d", id)
}

//	func publishOrderCreated(client *pubsub.Client, order OrderCreated) error {
//		ctx := context.Background()
//		topicID := os.Getenv("TOPIC_ID")
//		topic := client.Topic(topicID)
//
//		// publish order created event
//		result := topic.Publish(ctx, &pubsub.Message{
//			Data: []byte("order created"),
//		})
//
//		// block until publish is finished
//		_, err := result.Get(ctx)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
func createAndConfigureClient() (*pubsub.Client, error) {
	// get envs
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")

	if projectID == "" || topicID == "" {
		return nil, fmt.Errorf("PROJECT_ID and TOPIC_ID must be set")
	}

	// create client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	client.Topic(topicID)

	return client, nil
}

func handleOrder(ctx *gin.Context, client *pubsub.Client) {
	var order Order
	if err := ctx.BindJSON(&order); err != nil {
		ctx.JSON(400, gin.H{"message": "invalid request"})
		return
	}
	// generate order id
	id := generateOrderID()

	o := OrderCreated{
		Order: order,
		Id:    id,
	}

	println(o.Id)
	println(o.Order.Email)
	println(client)

	// todo order = c.Request.Bodycreate order in inmem db doesn't matter

	// publish order created event
	//err := publishOrderCreated(client, o)

	//if err != nil {
	//	c.JSON(500, gin.H{"message": "internal server error"})
	//	return
	//}

	ctx.JSON(200, gin.H{"message": "order created"})
}

// ApiMiddleware add client to context
func ApiMiddleware(client *pubsub.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("client", client)
		ctx.Next()
	}
}

func main() {
	client, err := createAndConfigureClient()
	if err != nil {
		// todo better exit
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(ApiMiddleware(client))
	r.POST("/order", func(ctx *gin.Context) {
		client := ctx.MustGet("client").(*pubsub.Client)
		handleOrder(ctx, client)
	})
	log.Fatal(r.Run(":8080"))
}
