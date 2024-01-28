package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// Order represents an order from the "frontend" or some other service doesn't matter
type Order struct {
	Email   string `json:"email" binding:"required"`
	Product string `json:"product" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
}

// OrderCreated represents an order created event, so the order created from the order received
type OrderCreated struct {
	Order Order  `json:"order"`
	Id    string `json:"id"`
}

// generateOrderID generates a random order id just because why not
func generateOrderID() string {
	id := rand.Intn(99999)
	return fmt.Sprintf("%05d", id)
}

// publishOrderCreated publishes an order created event to pubsub
func publishOrderCreated(client *pubsub.Client, order OrderCreated) error {
	ctx := context.Background()
	topicID := os.Getenv("TOPIC_ID")
	topic := client.Topic(topicID)

	// publish order created event
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("order created"),
	})

	// block until publish is finished
	msgID, err := result.Get(ctx)
	if err != nil {
		return err
	}
	log.Printf("published order created event with id %v", msgID)

	return nil
}

// createAndConfigureClient creates a pubsub client and configures it with projectid and topic
func createAndConfigureClient() (*pubsub.Client, error) {
	// get envs
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")

	if projectID == "" || topicID == "" {
		return nil, fmt.Errorf("PROJECT_ID and TOPIC_ID must be set")
	}

	// create client
	ctx := context.Background()
	authjson := os.Getenv("AUTH_JSON")
	opts := option.WithCredentialsFile(authjson)
	client, err := pubsub.NewClient(ctx, projectID, opts)

	if err != nil {
		return nil, err
	}
	client.Topic(topicID)

	return client, nil
}

// handleOrder handles an order request route
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

	// publish order created event
	err := publishOrderCreated(client, o)

	if err != nil {
		ctx.JSON(500, gin.H{"message": "internal server error"})
		log.Printf("error publishing order created event: %v", err)
		return
	}

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
	// create client
	client, err := createAndConfigureClient()
	if err != nil {
		// todo better exit
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// add client to context
	r.Use(ApiMiddleware(client))
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "ok"})
	})
	r.POST("/order", func(ctx *gin.Context) {
		// get client from context
		client := ctx.MustGet("client").(*pubsub.Client)
		handleOrder(ctx, client)
	})
	log.Fatal(r.Run(":8080"))
}
