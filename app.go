package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	ginTemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
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
	a.Engine.HTMLRender = ginTemplate.Default()
	a.initializeRoutes()
}
func (a *App) Run(port string) {
	a.Engine.Run(port)
}
