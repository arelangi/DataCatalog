package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Lineage struct {
	DerivedFrom string `json:"derived_from"`
	DatasetID   int64  `json:"dataset_id"`
}

func (a *App) registerLineageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errResp ErrorResponse
		var lineageRequest Lineage
		if err := c.Bind(&lineageRequest); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		fmt.Println(lineageRequest)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
