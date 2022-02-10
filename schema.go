package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *App) registerSchemaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var schemaRequest SchemaRequest
		var errResp ErrorResponse

		//Extract the request body
		schemaRequest.CatalogDatasetID, _ = strconv.ParseInt(c.PostForm("DatasetID"), 10, 64)
		schemaRequest.CatalogDatasetUUID, _ = strconv.ParseInt(c.PostForm("DatasetUUID"), 10, 64)

		//Extract the file
		file, err := c.FormFile("file")
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
		}
		if strings.Contains(file.Filename, ".avsc") {
			schemaRequest.SchemaType = "AVRO"
		} else if strings.Contains(file.Filename, ".proto") {
			schemaRequest.SchemaType = "PROTOBUF"
		}
		content, err := file.Open()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusPartialContent, ErrorResponse{Error: err.Error(), Message: "Failed to open file"})
			return
		}
		body, err := ioutil.ReadAll(content)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusPartialContent, ErrorResponse{Error: err.Error(), Message: "Failed to read file content"})
			return
		}
		schemaRequest.SchemaBody = string(body)

		//Construct the response
		response := a.constructRegistryResponseObject(&schemaRequest, body)

		//Register the schema
		err = a.registerSchema(&schemaRequest)
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = "Failed to register the schema"
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		//Generate and return the success response and the URL at which to publish the data to go to Kafka
		c.JSON(http.StatusOK, response)
	}
}

//Construct the response object for the REST call to register schema
func (a *App) constructRegistryResponseObject(schemaRequest *SchemaRequest, fileContent []byte) (response SchemaResponse) {
	var err error
	var fields AvroExtractResponse
	var failureResponse ErrorResponse

	response.DatasetID = schemaRequest.CatalogDatasetID
	response.Headers = []string{"--header 'Accept: application/vnd.kafka.v2+json, application/vnd.kafka+json, application/json'", "--header 'Content-Type: application/vnd.kafka.avro.v2+json'"}
	response.URL = fmt.Sprintf("http://restproxy/9082/topics/%s", schemaRequest.DatasetName)
	/*http://schemaregistry:8082/subjects/user-value/versions/latest/schema*/
	response.SchemaRegistryURL = fmt.Sprintf("http://schemaregistry:8082/schemas/%s-value/versions/latest/schema", schemaRequest.DatasetName)

	//Send this to the content to the AvroExtractor to extract the fields
	avroExtractorURL := "http://localhost:2602/demo/getFileContent"
	request := AvroExtractRequest{Schema: string(fileContent)}

	headersMap := map[string]string{
		"Content-Type": "application/json",
	}
	err = makePostCall(request, headersMap, avroExtractorURL, &fields, &failureResponse)
	if err != nil {
		log.Println("Failed to extract schema from the avro file ")
	}
	response.Fields = fields
	schemaRequest.Fields = response.Fields.Fields
	return
}

/*Once we have the file content, we need to do a few things.
0. Get the datsetname so the corrsponding kafka topic can be created and registered
1. Send this out to the schema registry to see if it's a valid request
2. If and when it is registered, update the status in the metadata table to schema applied
3. Insert record into the schemaregistry table with the corresponding ids
*/
func (a *App) registerSchema(schemaRequest *SchemaRequest) (err error) {
	kafkaRequestObj := KafkaSchemaRegistryRequest{Schema: schemaRequest.SchemaBody, SchemaType: schemaRequest.SchemaType}
	//Get the datasetname which will the Kafka topic name
	_, kafkaRequestObj.DatasetName = a.getDatasetSetName(schemaRequest.CatalogDatasetID)
	if len(kafkaRequestObj.DatasetName) == 0 {
		err = errors.New(fmt.Sprintf("Failed to find a dataset with the dataset id: %d", schemaRequest.CatalogDatasetID))
		return
	}
	schemaRequest.DatasetName = kafkaRequestObj.DatasetName

	//Make request to schema registry to register the schema
	schemaRegistryURL := fmt.Sprintf("http://schemaregistry:8082/subjects/%s-value/versions", schemaRequest.DatasetName)
	var failureResponse SchemaRegistrySuccessResponse
	headersMap := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, application/json",
	}

	err = makePostCall(kafkaRequestObj, headersMap, schemaRegistryURL, &schemaRequest, &failureResponse)
	if err != nil {
		log.Println("Failed to register the schema")
		return
	}
	//Insert the schemaregistry id to DB and Update metadata  status in the DB
	a.saveSchemaReferenceToDB(schemaRequest, "schema_applied")

	return
}

func (a *App) getDatasetSetName(id int64) (err error, datasetName string) {
	return a.DB.QueryRow("SELECT dataset_name FROM datacatalog.public.metadata where dataset_id=$1", id).Scan(&datasetName), datasetName
}

//Save the schema information to db and update the metadata status in the DB
func (a *App) saveSchemaReferenceToDB(schemaRequest *SchemaRequest, status string) (err error) {
	//Begin transaction to update the database
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	//Insert the kafka id into the schemaregistrymapping table
	_, err = tx.Exec("insert into datacatalog.public.schemaregistrymapping(kafka_registry_schema_id, dataset_id) VALUES($1, $2) on conflict do nothing", schemaRequest.KafkaRegistryID, schemaRequest.CatalogDatasetID)
	if err != nil {
		return err
	}

	for k, v := range schemaRequest.Fields {
		_, err = tx.Exec("insert into datacatalog.public.fields(dataset_id, field_id, name, description, types) VALUES($1, $2, $3, $4, $5) on conflict do nothing", schemaRequest.CatalogDatasetID, k, v.Name, v.Doc, v.Type)
		if err != nil {
			return err
		}
	}

	//Update the metadata tier status in the metadata table
	_, err = tx.Exec("UPDATE datacatalog.public.metadata set metadata_status=$1 where dataset_id=$2", status, schemaRequest.CatalogDatasetID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		//To-do: De-register from schema registry
		return err
	}

	return
}
