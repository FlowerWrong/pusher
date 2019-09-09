package main

import (
	"errors"
	"flag"
	"math/rand"
	"runtime"
	"time"

	"github.com/FlowerWrong/pusher"
	"github.com/FlowerWrong/pusher/api"
	"github.com/FlowerWrong/pusher/env"
	"github.com/FlowerWrong/pusher/log"
	"github.com/FlowerWrong/pusher/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	configFile := flag.String("config", "./settings.yml", "config file path")
	flag.Parse()
	err := pusher.Setup(*configFile)
	if err != nil {
		log.Panic(err)
	}

	if !pusher.ValidAppID(viper.GetString("APP_ID")) {
		log.Panic(errors.New("Invalid app id"))
	}

	log.Infoln("Pusher launch in", pusher.AppEnv)

	hub := pusher.NewHub()
	go hub.Run()

	router := gin.New()
	router.Use(middlewares.Logger(log.Logger()))
	router.Use(gin.Recovery())

	signatureGroup := router.Group("/apps", middlewares.Signature())
	{
		signatureGroup.POST("/:app_id/events", api.EventTrigger)
		signatureGroup.POST("/:app_id/batch_events", api.BatchEventTrigger)
		signatureGroup.GET("/:app_id/channels", api.ChannelIndex)
		signatureGroup.GET("/:app_id/channels/:channel_name", api.ChannelShow)
		signatureGroup.GET("/:app_id/channels/:channel_name/users", api.ChannelUsers)
	}

	// @doc https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol#websocket-connection
	// eg: ws://ws-ap1.pusher.com:80/app/APP_KEY?client=js&version=4.4&protocol=5
	router.GET("/app/:key", func(c *gin.Context) {
		appKey := c.Param("key")
		client := c.Query("client")
		version := c.Query("version")
		protocol := c.Query("protocol")
		pusher.ServeWs(hub, c.Writer, c.Request, appKey, client, version, protocol)
	})

	_ = router.Run(":" + env.Get("PORT", "8100"))
}
