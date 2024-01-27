package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pubsub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func Test_handOrder(t *testing.T) {
	tests := map[string]struct {
		wantCode int
		request  *http.Request
		w        *httptest.ResponseRecorder
		context  *gin.Context
	}{
		"valid request": {
			wantCode: 200,
			request:  httptest.NewRequest("POST", "/order", strings.NewReader(`{"email": "email.com", "product": "car", "amount": 99}`)),
			w:        httptest.NewRecorder(),
			context:  &gin.Context{},
		},
		"invalid request: missing email": {
			wantCode: 400,
			request:  httptest.NewRequest("POST", "/order", strings.NewReader(`{"product": "car", "amount": 99}`)),
			w:        httptest.NewRecorder(),
			context:  &gin.Context{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := GetTestContext(tc.w, tc.request)
			client, _ := pubsub.NewClient(ctx, "go-pubsub")
			handleOrder(ctx, client)
			if tc.w.Code != tc.wantCode {
				t.Fatalf("got %v, want %v", tc.w.Code, tc.wantCode)
			}

		})
	}
}

func TestWithMock(t *testing.T) {

	ctx := GetTestContext(httptest.NewRecorder(), httptest.NewRequest("POST", "/order", strings.NewReader(`{"email": "email.com", "product": "car", "amount": 99}`)))
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, err := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	handleOrder(ctx, client)
	if ctx.Writer.Status() != 200 {
		t.Fatalf("got %v, want %v", ctx.Writer.Status(), 200)
	}
}

func GetTestContext(w *httptest.ResponseRecorder, r *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c
}
