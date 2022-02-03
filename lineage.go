package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Lineage struct {
	DerivedFrom string    `json:"derived_from"`
	DatasetID   int64     `json:"dataset_id"`
	DatasetUUID uuid.UUID `json:"dataset_uuid"`
	//Audit columns
	CreatedDate     time.Time `json:"-"`
	LastUpdatedTime time.Time `json:"-"`
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

		if err := lineageRequest.createLineageRecord(a); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Failed to create account")
			c.JSON(http.StatusNotFound, errResp)
			return
		}
		c.JSON(http.StatusOK, lineageRequest)
	}
}

func (m *Lineage) createLineageRecord(app *App) (err error) {
	_, err = app.DB.Exec("INSERT into datacatalog.public.lineage(dataset_id, dataset_uuid, derived_from) VALUES ($1, $2, $3)", m.DatasetID, m.DatasetUUID, m.DerivedFrom)

	_, err = app.DB.Exec("UPDATE datacatalog.public.metadata set metadata_status=$1 where dataset_id=$2", "lineage_applied", m.DatasetID)

	return err
}
