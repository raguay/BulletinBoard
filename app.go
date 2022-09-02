package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/url"
	"time"
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

type Dialog struct {
	Html   string `json:"html" binding:"required"`
	Width  int    `json:"width" binding:"required"`
	Height int    `json:"height" binding:"required"`
	X      int    `json:"x" binding:"required"`
	Y      int    `json:"y" binding:"required"`
}

type processData struct {
	c       *gin.Context
	ctx     context.Context
	running bool
}

func newProcessData() *processData {
	return &processData{}
}

func (p *processData) init(c *gin.Context, ctx context.Context) {
	p.c = c
	p.ctx = ctx
	p.running = true
}

func (p *processData) optionalData(msg string) {
	p.c.JSON(http.StatusOK, msg)
	rt.EventsOff(p.ctx, "dialogreturn")
	p.running = false
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
		retData := newProcessData()
		retData.init(c, ctx)
		//rt.EventsOn(ctx, "dialogreturn", retData)
		for retData.running {
			time.Sleep(30 * time.Second)
		}
	})

	//
	// Run the server.
	//
	r.Run(":9697")
}
