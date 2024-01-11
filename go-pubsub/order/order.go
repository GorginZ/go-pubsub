package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
)

// Order represents an order from the frontend
type Order struct {
	Email   string `json:"email"`
	Product string `json:"product"`
	Amount  int    `json:"amount"`
}

type OrderCreated struct {
	Order Order  `json:"order"`
	Id    string `json:"id"`
}

func generateOrderID() string {
	id := rand.Intn(99999)
	return fmt.Sprintf("%05d", id)
}
//
//func publishOrderCreated(client *pubsub.Client, order OrderCreated) error {
//	ctx := context.Background()
//	topicID := os.Getenv("TOPIC_ID")
//	topic := client.Topic(topicID)
//
//	// publish order created event
//	result := topic.Publish(ctx, &pubsub.Message{
//		Data: []byte("order created"),
//	})
//
//	// block until publish is finished
//	_, err := result.Get(ctx)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// handle order will generate a order id and create a new order in an inmemory database and publish the order to pubsub
//func handleOrder(c *gin.Context, client *pubsub.Client) {
//	var order Order
//	if err := c.BindJSON(&order); err != nil {
//		c.JSON(400, gin.H{"message": "invalid request"})
//		return
//	}
//	// generate order id
//	id := generateOrderID()
//
//	o := OrderCreated{
//		Order: order,
//		Id:    id,
//	}
//
//	// todo create order in inmem db doesn't matter
//
//	// publish order created event
////	err := publishOrderCreated(client, o)
//
////	if err != nil {
////		c.JSON(500, gin.H{"message": "internal server error"})
////		return
////	}
//
//	c.JSON(200, gin.H{"message": "order created"})
//}
//
//func createAndConfigureClient() (*pubsub.Client, error) {
//	// get envs
//	projectID := os.Getenv("PROJECT_ID")
//	topicID := os.Getenv("TOPIC_ID")
//
//	if projectID == "" || topicID == "" {
//		return nil, fmt.Errorf("PROJECT_ID and TOPIC_ID must be set")
//	}
//
//	// create client
//	ctx := context.Background()
//	client, err := pubsub.NewClient(ctx, projectID)
//	client.Topic(topicID)
//
//	return client, err
//}
func handleOrder(c *gin.Context) {
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}
	// generate order id
	id := generateOrderID()

	o := OrderCreated{
		Order: order,
		Id:    id,
	}

	println(o.Id)

	// todo order = c.Request.Bodycreate order in inmem db doesn't matter

	// publish order created event
	//err := publishOrderCreated(client, o)

	//if err != nil {
	//	c.JSON(500, gin.H{"message": "internal server error"})
	//	return
	//}

	c.JSON(200, gin.H{"message": "order created"})
}

func main() {
	//c, err := createAndConfigureClient()
	//if err != nil {
	//	// todo better exit
	//	log.Fatal(err)
	//}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/order", handleOrder)
	log.Fatal(r.Run(":8080"))
}
