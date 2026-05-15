package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	api "github.com/WeatherGod3218/iat-animals-rewrite/internal/api"
)

const PORT string = "3000"

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./public")

	router.Use(cors.Default())

	router.GET("/", api.GetHomepage)
	router.GET("/get-data", api.GetData)

	router.Run(":" + PORT)
}
