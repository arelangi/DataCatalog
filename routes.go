package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func (a *App) initializeRoutes() {
	log.Println("In initializeRoutes")

	// Swagger Documention
	a.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Static content
	a.Engine.Use(static.Serve("/", static.LocalFile("./views", true)))

	//Group register related routes together
	registerRoutes := a.Engine.Group("/register")
	{

		//Page rendering
		//Render the register page
		registerRoutes.GET("/start", a.showRegisterPage())
		registerRoutes.GET("/dataclassification/:dataset_id", a.showDataClassificationPage())

		//API Calls

		//Handle the POST requests at /register/metadata
		registerRoutes.POST("/metadata", a.registerMetadataHandler())

		//Handle the POST requests at /register/schema
		registerRoutes.POST("/schema", a.registerSchemaHandler())

		//Handle the POST requests at /register/lineage
		registerRoutes.POST("/lineage", a.registerLineageHandler())
	}

	//Test
	a.Engine.GET("/security", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "security.html", gin.H{
			"title": "Security file title!!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

}
