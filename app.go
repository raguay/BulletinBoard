package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/url"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) domReady(ctx context.Context) {

}

func (a *App) shutdown(ctx context.Context) {

}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	//
	// We need to start the backend and setup the signaling.
	//
	go backend(ctx)
}

type Msg struct {
	Message string `json:"msg" xml:"user"  binding:"required"`
}

func backend(ctx context.Context) {
	//
	// This will have the web server backend for BulletinBoard.
	//
	r := gin.Default()
	r.Use(gin.Recovery())

	//
	// Define the message route. The message is given on the URI string and in the body.
	//
	r.GET("/api/message/:message", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "okay",
		})
		var json Msg
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		message := c.Param("message")
		messageBody := json.Message
		if messageBody != message {
			message = messageBody
		}

		message, err := url.QueryUnescape(message)
		if err != nil {
			// An error in decoding.
			message = ""
		}

		//
		// Send it to the frontend.
		//
		rt.EventsEmit(ctx, "message", message)
	})

	//
	// Run the server.
	//
	r.Run(":9697")
}
