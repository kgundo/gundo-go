package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kgundo/gundo-go/config"
	"github.com/kgundo/gundo-go/routes"
)

func main() {
	router := gin.New()
	config.Connect()
	routes.UserRoute(router)
	router.Run(":8080")
}
