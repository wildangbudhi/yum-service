package main

import (
	"log"

	"github.com/gin-gonic/gin"
	v1 "github.com/wildangbudhi/yum-service/depedencyinjection/v1"
	"github.com/wildangbudhi/yum-service/utils"
)

func main() {
	server, err := utils.NewServer()

	if err != nil {
		log.Fatal(err)
	}

	HealthCheckHandler(server)
	depedencyInjection(server)
	server.Router.Run(server.Config.ServerHost)

	err = server.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func depedencyInjection(server *utils.Server) {
	v1.AuthHTTPRestDI(server)
}

func HealthCheckHandler(server *utils.Server) {
	server.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
