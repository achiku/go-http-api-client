package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/achiku/testsvr"
)

func TestClient_Hello(t *testing.T) {
	ts := httptest.NewServer(testsvr.NewMux(DefaultHandlerMap, t))
	defer ts.Close()

	client := NewClient(TestNewConfig(ts.URL), &http.Client{}, nil)
	ctx := context.Background()
	req := &HelloRequest{
		Name: "achiku",
	}
	res, err := client.Hello(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != SuccessStatusCode {
		t.Errorf("want %d got %d", SuccessStatusCode, res.StatusCode)
	}
	t.Logf("%+v", res)
}
