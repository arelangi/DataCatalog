package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *App) RegisterDatasetToDatahub(id int64) (err error) {
	var ds DatasetSnapshotType
	var failureMessage DatahubAPIFailureResponse
	var successMessage struct{}

	datahubRestURL := "http://datahub:8080/entities?action=ingest"
	headersMap := map[string]string{}

	dataset, err := a.getCompleteDatasetByID(id)
	if err != nil {
		panic(err)
	}

	var metadata MetadataRequest
	a.getMetadata(id, &metadata)
	metadata.DatasetID = id

	//Build the urn for the object
	platform := "kafka"
	namespace := "PROD"
	ds.Urn = fmt.Sprintf("urn:li:dataset:(urn:li:dataPlatform:%s,%s,%s)", platform, metadata.DatasetName, namespace)

	//Get the owner aspect
	owner := metadata.getOwnerAspect()
	aspectOwnership := AspectType{Ownership: &owner}

	//Get the institutionalMemory aspect
	memory := metadata.getInstitutionalMemoryAspect()
	aspectInstitutionalMemory := AspectType{InstitutionalMemory: &memory}

	//Get the UpstreamLineage aspect
	if strings.ToLower(metadata.PlatformName) != "kafka" {
		//Then we have some upstream lineage. For now, we skip this section
		upstream := metadata.getUpstreamLineageAspect()
		aspectUpstreamLineage := AspectType{UpstreamLineage: &upstream}
		ds.Aspects = append(ds.Aspects, aspectUpstreamLineage)
	}

	//Get the DatasetProperties aspect
	properties := metadata.getDatasetPropertiesAspect()
	aspectDatasetProperties := AspectType{DatasetProperties: &properties}

	//Get the SchemaMetadata aspect
	schemaMetadata := metadata.getSchemaMetadataAspect(a)
	aspectSchemaMetaData := AspectType{SchemaMetaData: &schemaMetadata}

	//The following is where we build the request object and make the post call to datahub API
	ds.Aspects = append(ds.Aspects, aspectOwnership, aspectInstitutionalMemory, aspectDatasetProperties, aspectSchemaMetaData)
	req := DataHubRequestType{
		Entity: DataHubEntityType{
			Value: DataHubEntityValueType{
				DatasetSnapshot: ds,
			},
		},
	}

	dhReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal the request to JSON", err)
		return
	}
	fmt.Println("The following is the request we are sending to datahub")
	fmt.Println("------------------------------------------------------")
	fmt.Println(string(dhReq))
	fmt.Println("------------------------------------------------------")

	makePostCall(req, headersMap, datahubRestURL, &successMessage, &failureMessage)

	fmt.Println(req, ds, dataset)

	return err
}

func (a *App) registerDatahubHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errResp ErrorResponse

		datasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64)
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return

		}

		err = a.RegisterDatasetToDatahub(datasetID)
		if err != nil {
			fmt.Println("Failed to register to Datahub:", err)
			panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": err.Error()})
			return
		}

	}
}
