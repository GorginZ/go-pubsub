package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/pstest"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// HappyPublishReactor will give us back a message ID to fake out success of the google pubsub api call
type HappyPublishReactor struct{}

func (r *HappyPublishReactor) React(_ interface{}) (handled bool, ret interface{}, err error) {
	pbr := &pubsubpb.PublishResponse{MessageIds: []string{"61"}}

	return true, pbr, nil
}

func Test_handleOrder(t *testing.T) {
	tests := map[string]struct {
		wantCode int
		request  *http.Request
		w        *httptest.ResponseRecorder
		context  *gin.Context
		reactor  pstest.Reactor //  https://pkg.go.dev/cloud.google.com/go/pubsub/pstest?utm_source=godoc#Reactor
	}{
		"valid request": {
			wantCode: 200,
			request:  httptest.NewRequest("POST", "/order", strings.NewReader(`{"email": "email.com", "product": "car", "amount": 99}`)),
			w:        httptest.NewRecorder(),
			context:  &gin.Context{},
			reactor:  &HappyPublishReactor{}, //will return a msg ID
		},
		"invalid request: missing email": {
			wantCode: 400,
			request:  httptest.NewRequest("POST", "/order", strings.NewReader(`{"product": "car", "amount": 99}`)),
			w:        httptest.NewRecorder(),
			context:  &gin.Context{},
			reactor:  &HappyPublishReactor{}, //not really relevant this will "fail fast"
		},
		"pubsub error": {
			wantCode: 500,
			request:  httptest.NewRequest("POST", "/order", strings.NewReader(`{"email": "email.com", "product": "car", "amount": 99}`)),
			w:        httptest.NewRecorder(),
			context:  &gin.Context{},
			// https://github.com/googleapis/google-cloud-go/blob/pubsub/v1.36.0/pubsub/pstest/fake.go#L1470
			reactor: pstest.WithErrorInjection("Publish", 200, "georgia error").Reactor, //will give an error
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := GetTestContext(tc.w, tc.request)

			//set up server
			opt := pstest.ServerReactorOption{
				FuncName: "Publish",
				Reactor:  tc.reactor,
			}
			srv := pstest.NewServer(opt)
			// srv := pstest.NewServer()
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
			if tc.w.Code != tc.wantCode {
				t.Fatalf("got %v, want %v", tc.w.Code, tc.wantCode)
			}

		})
	}
}

func GetTestContext(w *httptest.ResponseRecorder, r *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c
}
