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

func getAirtableData(c *gin.Context) ([]airtable.AirtableRecord, error) {
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

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetData(c *gin.Context) {
	data, err := getAirtableData(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "chooseAirtableMiddleware"}).Fatal("error fetching airtable!")

	}

	images := make([]string, 0)
	test_stimuli := make([]airtable.AirtableClientResponse, 0)
	categoryWordImage := "words_animate_Cat1.png"
	categoryDisplay := map[int][][]string{
		1: make([][]string, 2),
		2: make([][]string, 2),
		3: make([][]string, 2),
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Fields.Trial < data[j].Fields.Trial
	})

	for _, record := range data {

		if record.Fields.Stimulus == "inert" && record.Fields.CorrectKey == "d" {
			categoryWordImage = "words_animate_Cat2.png"
		}

		if record.Fields.StimulusType == "image" {
			images = append(images, record.Fields.Stimulus)
		}

		var stimuliField airtable.AirtableClientResponse
		onLeft := record.Fields.CorrectKey == "d"
		if onLeft {
			stimuliField = airtable.GetResponse(record.Fields, "left")
		} else {
			stimuliField = airtable.GetResponse(record.Fields, "right")
		}
		test_stimuli = append(test_stimuli, stimuliField)

		idx := 1
		if onLeft {
			idx = 0
		}

		if !contains(categoryDisplay[record.Fields.Block][idx], record.Fields.Category) {
			categoryDisplay[record.Fields.Block][idx] = append(categoryDisplay[record.Fields.Block][idx], record.Fields.Category)
		}
	}

	for k := range categoryDisplay {
		sort.Slice(categoryDisplay[k][0], func(i, j int) bool {
			return len(categoryDisplay[k][0][i]) < len(categoryDisplay[k][0][j])
		})

		sort.Slice(categoryDisplay[k][1], func(i, j int) bool {
			return len(categoryDisplay[k][1][i]) < len(categoryDisplay[k][1][j])
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"test_stimuli":        test_stimuli,
		"images":              images,
		"category_display":    categoryDisplay,
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
