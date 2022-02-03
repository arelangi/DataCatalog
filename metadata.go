package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetadataRequest struct {
	//Basic Attributes
	DatasetName                   string `json:"dataset_name"`
	DatasetLogicalName            string `json:"dataset_logical_name"`
	DatasetDescription            string `json:"dataset_description"`
	DatasetType                   string `json:"dataset_type"`
	DatasetSource                 string `json:"dataset_source"`
	DatasetShare                  string `json:"dataset_share"`
	DatasetRetention              int64  `json:"dataset_retention"`
	DatasetRetentionJustification string `json:"dataset_retention_justification"`
	DatasetArrivalFrequency       string `json:"dataset_arrival_frequency"`

	//Ownership Attributes
	Organization string `json:"organization"`
	Product      string `json:"product"`
	Team         string `json:"team"`
	DataSteward  string `json:"data_steward"`

	//Platform attributes
	PlatformName string `json:"platform_name"`
	HostName     string `json:"host_name"`
	DatabaseName string `json:"database_name"`
	SchemaName   string `json:"schema_name"`

	//Security & Privacy attributes
	DataClassiffication string `json:"data_classiffication"`

	//IDs
	DatasetID   int64     `json:"dataset_id"`
	DatasetUUID uuid.UUID `json:"dataset_uuid"`

	//Status attributes
	MetadataStatus string `json:"metadata_status"`

	//Audit columns
	CreatedDate     time.Time `json:"created_date"`
	LastUpdatedTime time.Time `json:"last_updated_time"`
}

func (a *App) registerMetadataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errResp ErrorResponse
		var metadataRequest MetadataRequest
		if err := c.Bind(&metadataRequest); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		fmt.Println(metadataRequest)

		if err := metadataRequest.createMetadataRecord(a); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Failed to create account")
			c.JSON(http.StatusNotFound, errResp)
			return
		}
		c.JSON(http.StatusOK, metadataRequest)
	}
}

func (m *MetadataRequest) createMetadataRecord(app *App) (err error) {
	err = app.DB.QueryRow("INSERT INTO datacatalog.public.metadata(dataset_name,  dataset_logical_name,  dataset_description,  dataset_type,  dataset_source,  dataset_share,  dataset_retention,  dataset_retention_justification,  dataset_arrival_frequency, organization,product,team,data_steward,platform_name,host_name,database_name,schema_name,data_classiffication) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) returning dataset_id, dataset_uuid, metadata_status", m.DatasetName, m.DatasetLogicalName, m.DatasetDescription, m.DatasetType, m.DatasetSource, m.DatasetShare, m.DatasetRetention, m.DatasetRetentionJustification, m.DatasetArrivalFrequency, m.Organization, m.Product, m.Team, m.DataSteward, m.PlatformName, m.HostName, m.DatabaseName, m.SchemaName, m.DataClassiffication).Scan(&m.DatasetID, &m.DatasetUUID, &m.MetadataStatus)

	return err
}
