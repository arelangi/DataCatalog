package main

import (
	"log"

	"github.com/gin-gonic/contrib/static"
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
		registerRoutes.GET("/classifydata/:dataset_id", a.showDataClassificationPage())
		registerRoutes.GET("/dataquality/:dataset_id", a.showDataQualityPage())

		//API Calls
		//Handle the POST requests at /register/metadata
		registerRoutes.POST("/metadata", a.registerMetadataHandler())

		//Handle the POST requests at /register/schema
		registerRoutes.POST("/schema", a.registerSchemaHandler())

		//Handle the POST request at /register/classification
		registerRoutes.POST("/classification", a.dataClassificationHandler())

		//Handle the POST request at /register/quality
		registerRoutes.POST("/quality", a.dataQualityHandler())

		//Handle the POST requests at /register/lineage
		registerRoutes.POST("/lineage", a.registerLineageHandler())
	}

	dataStewardRoutes := a.Engine.Group("/ds")
	{
		//Page rendering
		dataStewardRoutes.GET("/review/:dataset_id", a.showReviewPage())

		//Handle the data steward approval
		dataStewardRoutes.GET("/approval/:dataset_id", a.approveDatasetHandler())
	}

}
