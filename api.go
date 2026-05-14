package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetData() {
	router := gin.Default()

	router.Use(cors.Default())

	router.Run(":" + PORT)
}

func replaceNullWithString(obj map[string]interface{}) {
	for key, value := range obj {
		if value == nil {
			obj[key] = "null"
		} else if nested, ok := value.(map[string]interface{}); ok {
			replaceNullWithString(nested)
		}
	}
}

func SubmitResults(c *gin.Context) {
	var results map[string]interface{}

	err := c.ShouldBindJSON(&results)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	replaceNullWithString(results)

}
