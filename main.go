package main

import (
	"embed"
	"html/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	api "github.com/WeatherGod3218/iat-animals-rewrite/internal/api"
)

const PORT string = "3000"

var templateFS embed.FS

func templateFromEmbed() *template.Template {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*"))
	return tmpl
}

func main() {

	router := gin.Default()
	router.SetHTMLTemplate(templateFromEmbed())
	router.Static("/static", "./public")

	router.Use(cors.Default())

	router.GET("/", api.GetHomepage)
	router.GET("/get-data", api.GetData)

	router.Run(":" + PORT)
}
