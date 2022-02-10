package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) dataClassificationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errResp ErrorResponse
		fmt.Println(errResp)

		f, err := c.MultipartForm()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		response, err := a.saveClassificationToDB(f.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.JSON(http.StatusOK, response)
	}
}

func (a *App) saveClassificationToDB(data map[string][]string) (retVal map[string]string, err error) {
	retVal = make(map[string]string)
	datasetID := data["datasetid"][0]
	delete(data, "datasetid")
	for k, v := range data {
		_, err = a.DB.Exec("UPDATE datacatalog.public.fields SET classification=$1 where dataset_id=$2 and field_id=$3", v[0], datasetID, k)
		if err != nil {
			return
		}
		retVal[k] = v[0]
	}
	return
}
