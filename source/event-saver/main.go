package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	dapr "github.com/mchmarny/godapr/v1"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// dapr
	daprClient Client

	// service
	servicePort = env.MustGetEnvVar("PORT", "8080")
	sourceTopic = env.MustGetEnvVar("SOURCE_TOPIC_NAME", "events")
	targetStore = env.MustGetEnvVar("TARGET_STORE_NAME", "store")
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	daprClient = dapr.NewClient()

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// simple routes
	r.GET("/", rootHandler)
	r.GET("/dapr/subscribe", subscribeHandler)

	// topic handler route
	eventHandlerRoute := fmt.Sprintf("/%s", sourceTopic)
	logger.Printf("viewer route: %s", eventHandlerRoute)
	r.POST(eventHandlerRoute, eventHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := http.ListenAndServe(hostPort, &ochttp.Handler{Handler: r}); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

// Options midleware
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Client is the minim client support for testing
type Client interface {
	SaveState(ctx trace.SpanContext, store, key string, data interface{}) error
}
