package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestEventHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	daprClient = GetTestClient()

	r := gin.Default()
	r.POST("/", subscribeHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./event.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestTopicListHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	daprClient = GetTestClient()

	r := gin.Default()
	r.GET("/", subscribeHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var out []string
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Len(t, out, 1)

}

func GetTestClient() *TestClient {
	return &TestClient{}
}

var (
	// test test client against local interace
	_ = Client(&TestClient{})
)

type TestClient struct {
}

func (c *TestClient) SaveState(ctx trace.SpanContext, store, key string, data interface{}) error {
	return nil
}
