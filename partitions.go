package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) registerPartitionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var partitionRequest ParitionRequest
		var errResp ErrorResponse

		if err := c.Bind(&partitionRequest); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		fmt.Println("This is what I received ", partitionRequest)

		if err := a.UpdatePartitionFields(partitionRequest); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

func (a *App) UpdatePartitionFields(req ParitionRequest) (err error) {
	//Update the fields tables with these details

	fmt.Println("Saving the result ", req)
	for _, v := range req.PrimaryKeys {
		_, err = a.DB.Exec("update datacatalog.public.fields set primarykeyfield=true where dataset_id=$1 and field_id=$2", req.DatasetID, v)
		if err != nil {
			return
		}
	}

	_, err = a.DB.Exec("update datacatalog.public.fields set partitionfield=true where dataset_id=$1 and field_id=$2", req.DatasetID, req.PartitionPath)
	return
}

type ParitionRequest struct {
	PartitionPath int64   `json:"partition_path"`
	PrimaryKeys   []int64 `json:"primary_keys"`
	DatasetID     int64   `json:"dataset_id"`
}
