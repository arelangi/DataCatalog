package main

import "github.com/gin-gonic/gin"

func (a *App) showRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		render(c, gin.H{
			"title": "Register Data",
		}, "registerdata.html")
	}
}

func (a *App) showDataClassificationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
