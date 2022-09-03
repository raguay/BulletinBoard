package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/url"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
	srv *http.Server
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) domReady(ctx context.Context) {

}

func (a *App) shutdown(ctx context.Context) {
	a.srv.Shutdown(ctx)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	//
	// We need to start the backend and setup the signaling.
	//
	go backend(a, ctx)
}

type Msg struct {
	Message string `json:"msg" xml:"user"  binding:"required"`
}

type Dialog struct {
	Html   string `json:"html" binding:"required"`
	Width  int    `json:"width" binding:"required"`
	Height int    `json:"height" binding:"required"`
	X      int    `json:"x" binding:"required"`
	Y      int    `json:"y" binding:"required"`
}

func backend(a *App, ctx context.Context) {
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
	// Define the append to message route. This one is URI only.
	//
	r.GET("/api/message/append/:message", func(c *gin.Context) {
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
		rt.EventsEmit(ctx, "append", message)
	})

	//
	// Add the dialog route for user defined raw dialogs.
	//
	r.PUT("/api/dialog", func(c *gin.Context) {
		var json Dialog
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//
		// Send it to the frontend.
		//
		rt.EventsEmit(ctx, "dialog", json)

		//
		// Get the return.
		//
		running := true
		rt.EventsOn(ctx, "dialogreturn", func(optionalData ...interface{}) {
			c.JSON(http.StatusOK, optionalData)
			running = false
			rt.EventsOff(ctx, "dialogreturn")
		})
		for running {
			time.Sleep(time.Second)
		}
	})

	//
	// Add the quit route.
	//
	r.GET("/api/quit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "okay",
		})

		//
		// Exit the application.
		//
		os.Exit(0)
	})

	//
	// Run the server.
	//
	r.Run(":9697")
}
