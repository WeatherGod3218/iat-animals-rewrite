package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/iat-animals-rewrite/internal/airtable"
	"github.com/WeatherGod3218/iat-animals-rewrite/internal/firebase"

	"github.com/WeatherGod3218/iat-animals-rewrite/logging"
	"github.com/sirupsen/logrus"
)

func replaceNullWithString(obj map[string]interface{}) {
	for key, value := range obj {
		if value == nil {
			obj[key] = "null"
		} else if nested, ok := value.(map[string]interface{}); ok {
			replaceNullWithString(nested)
		}
	}
}

func chooseAirtableMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName, err := firebase.GetLowestTable(c)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "chooseAirtableMiddleware"}).Fatal("error deciding which airtable to use!")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error selecting dataset to display.", "error": err.Error()})
			return
		}

		tableURI := os.Getenv("AIRTABLE_TABLE" + tableName)

		airTable, err := airtable.GetAirtableURI(tableURI)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "chooseAirtableMiddleware"}).Fatal("error fetching airtable!")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching dataset.", "error": err.Error()})
			return
		}

		c.Set("airtable", airTable)
		c.Next()
	}
}

func GetHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func GetData(c *gin.Context) {

}

func SubmitResults(c *gin.Context) {
	var results map[string]interface{}

	err := c.ShouldBindJSON(&results)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	replaceNullWithString(results)

	err = firebase.PushToDatabase(c, results)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error updating database!")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving results.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully saved results!"})
}
