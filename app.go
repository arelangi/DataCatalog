package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type App struct {
	Engine *gin.Engine
	DB     *sql.DB
}

func (a *App) Initialize(dbName, dbUser, dbPassword string) {

	//Database connection
	host := "postgresserver"
	port := 5432

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, dbName)
	fmt.Println(psqlconn)
	var err error //
	a.DB, err = sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal("Failed to connect", err.Error())
	}
	//Connectivity check
	err = a.DB.Ping()
	if err != nil {
		panic(err)
	}

	// Server initlization
	a.Engine = gin.Default()
	a.initializeRoutes()
}
func (a *App) Run(port string) {
	a.Engine.Run(port)
}

func (a *App) initializeRoutes() {
	fmt.Println("In initializeRoutes")
	// Swagger Documention
	a.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Static content
	a.Engine.Use(static.Serve("/", static.LocalFile("./views", true)))

	//API Calls
	a.Engine.POST("/registerMetadata", a.registerMetadataHandler())
	a.Engine.POST("/registerLineage", a.registerLineageHandler())

}
