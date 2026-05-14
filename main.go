package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/WeatherGod3218/iat-animals-rewrite/logging"
	"github.com/sirupsen/logrus"
)

const PORT string = "3000"

func main() {
	err := godotenv.Load()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error fetching .env file")
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(301, "/")
	})
	router.Run(":" + PORT)
}
