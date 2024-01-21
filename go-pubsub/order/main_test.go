package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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
			handleOrder(ctx)
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
