package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) showRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		render(c, gin.H{
			"title": "Register Data",
		}, "registerdata.html")
	}
}

func (a *App) showDataClassificationPage() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if the dataset ID is valid
		if catalogDatasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64); err == nil {
			// Check if the dataset exists
			if dataset, err := a.getDatasetByID(catalogDatasetID); err == nil {

				// Call the render function with the title, dataset and the name of the
				// template
				render(c, gin.H{
					"payload": dataset,
				}, "classifydata.html")

			} else {
				// If the dataset is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}

		} else {
			// If an invalid dataset ID is specified in the URL, abort with an error
			c.AbortWithStatus(http.StatusNotFound)
		}

	}
}

func (a *App) showDataQualityPage() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if the dataset ID is valid
		if catalogDatasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64); err == nil {
			// Check if the dataset exists
			if dataset, err := a.getDatasetByID(catalogDatasetID); err == nil {

				// Call the render function with the title, dataset and the name of the
				// template
				render(c, gin.H{
					"payload": dataset,
				}, "dataquality.html")

			} else {
				// If the dataset is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}

		} else {
			// If an invalid dataset ID is specified in the URL, abort with an error
			c.AbortWithStatus(http.StatusNotFound)
		}

	}
}

func (a *App) showApprovalPage() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
