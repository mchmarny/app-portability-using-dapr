package main

import (
	"fmt"
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

const (
	// SupportedCloudEventVersion indicates the version of CloudEvents suppored by this handler
	SupportedCloudEventVersion = "0.3"

	//SupportedCloudEventContentTye indicates the content type supported by this handlers
	SupportedCloudEventContentTye = "application/json"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"requestor": c.Request.Host,
		"version":   AppVersion,
	})
}

func subscribeHandler(c *gin.Context) {
	topics := []string{sourceTopic}
	logger.Printf("subscription tipics: %v", topics)
	c.JSON(http.StatusOK, topics)
}

func eventHandler(c *gin.Context) {
	httpFmt := tracecontext.HTTPFormat{}
	ctx, ok := httpFmt.SpanContextFromRequest(c.Request)
	if !ok {
		ctx = trace.SpanContext{}
	}

	logger.Printf("trace info: 0-%x-%x-%x",
		ctx.TraceID[:],
		ctx.SpanID[:],
		[]byte{byte(ctx.TraceOptions)})

	e := ce.NewEvent()
	if err := c.ShouldBindJSON(&e); err != nil {
		logger.Printf("error binding event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error processing your request, see logs for details",
		})
		return
	}

	// check the event version (0.3 supported)
	eventVersion := e.Context.GetSpecVersion()
	if eventVersion != SupportedCloudEventVersion {
		logger.Printf("invalid event spec version: %s", eventVersion)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": fmt.Sprintf("Invalid spec version (want: %s got: %s)",
				SupportedCloudEventVersion, eventVersion),
		})
		return
	}

	// check the content type (looking for application/json)
	eventContentType := e.Context.GetDataContentType()
	if eventContentType != SupportedCloudEventContentTye {
		logger.Printf("invalid event content type: %s", eventContentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": fmt.Sprintf("Invalid content type (want: %s got: %s)",
				SupportedCloudEventContentTye, eventContentType),
		})
		return
	}

	eventID := e.Context.GetID()
	logger.Printf("event ID: %s", eventID)

	eventContent := string(e.Data())
	logger.Printf("event content: %s", eventContent)

	err := daprClient.SaveState(ctx, targetStore, eventID, eventContent)
	if err != nil {
		logger.Printf("error saving state: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	c.JSON(http.StatusOK, nil)
}

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)
