package api

import (
	"net/http"
	"os"
	"sort"

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

func getAirtableData(c *gin.Context) ([]map[string]interface{}, error) {
	tableName, err := firebase.GetLowestTable(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error deciding which airtable to use!")
		return nil, err
	}

	tableURI := os.Getenv("AIRTABLE_TABLE" + tableName)

	airTable, err := airtable.GetAirtableURI(tableURI)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtab"}).Warn("error fetching airtable!")
		return nil, err
	}

	return airTable, nil
}

func GetHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func GetData(c *gin.Context) {
	data, err := getAirtableData(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "chooseAirtableMiddleware"}).Fatal("error fetching airtable!")

	}

	images := make([]string, 0)
	test_stimuli := make([]map[string]interface{}, 0)
	categoryWordImage := "words_animate_Cat1.png"

	sort.Slice(data, func(i, j int) bool {
		return data[i]["trial"].(float64) < data[j]["trial"].(float64)
	})

	for _, fields := range data {
		if fields["stimulus"] == "inert" && fields["correct_key"] == "d" {
			categoryWordImage = "words_animate_Cat2.png"
		}

		if fields["stimulus_type"] == "image" {
			if stimulus, ok := fields["stimulus"].(string); ok {
				images = append(images, stimulus)
			}
		}

		if fields["correct_key"] == "d" {
			fields["association"] = "left"
		} else {
			fields["association"] = "right"
		}
		test_stimuli = append(test_stimuli, fields)

	}

	c.JSON(http.StatusOK, gin.H{
		"test_stimuli":        "hi",
		"category_word_image": categoryWordImage,
	})
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
