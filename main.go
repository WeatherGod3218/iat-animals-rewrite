package main

import (
	"io/fs"
	"net/http"
	"os"

	"embed"
	"html/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	api "github.com/WeatherGod3218/iat-animals-rewrite/internal/api"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed templates/* public/*
var embeddedFS embed.FS

func templateFromEmbed() *template.Template {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*"))
	return tmpl
}

func main() {

	router := gin.Default()

	router.Use(cors.Default())

	tmpl := template.Must(template.ParseFS(embeddedFS, "templates/*"))
	router.SetHTMLTemplate(tmpl)

	publicFS, err := fs.Sub(embeddedFS, "public")
	if err != nil {
		panic(err)
	}

	router.StaticFS("/static", http.FS(publicFS))

	router.GET("/", api.GetHomepage)
	router.GET("/get-data", api.GetData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
