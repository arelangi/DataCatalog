package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SchemaRequest struct {
	CatalogDatasetID   int64  `json:"catalog_dataset_id"`
	CatalogDatasetUUID int64  `json:"catalog_dataset_uuid"`
	SchemaBody         string `json:"schema_body"`
	SchemaType         string `json:"schema_type"`
	DatasetName        string `json:"dataset_name"`

	KafkaRegistryID int64 `json:"kafka_registry_id"`
}

type SchemaResponse struct {
	SchemaRegistryURL string              `json:"schema_registry_url"`
	CurlCommand       string              `json:"curl_command"`
	URL               string              `json:"url"`
	Headers           []string            `json:"headers"`
	SchemaBody        string              `json:"schema_body"`
	Fields            AvroExtractResponse `json:"fields"`
}

//To-do: Parquet data format as well
type AvroRequest struct {
	ValueSchemaID int64    `json:"value_schema_id"`
	Records       []Record `json:"records"`
}

type Record struct {
	Value []string `json:"value"`
}

type Value struct {
	string
}

type KafkaSchemaRegistryRequest struct {
	DatasetName string `json:"dataset_name"`
	Schema      string `json:"schema"`
	SchemaType  string `json:"schemaType"`
}

//To-do: Interface for better error messages for kafka registry
type SchemaRegistrySuccessResponse struct {
	Id int64 `json:"id"`
}

type SchemaRegistrySuccessFailure struct {
	ErrorCode int64  `json:"error_code"`
	Message   string `json:"message"`
}

func (a *App) registerSchemaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		schemaRequest := SchemaRequest{}
		var errResp ErrorResponse
		file, err := c.FormFile("file")
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
		}

		schemaRequest.CatalogDatasetID, _ = strconv.ParseInt(c.PostForm("DatasetID"), 10, 64)
		schemaRequest.CatalogDatasetUUID, _ = strconv.ParseInt(c.PostForm("DatasetUUID"), 10, 64)

		if strings.Contains(file.Filename, ".avsc") {
			schemaRequest.SchemaType = "AVRO"
		} else if strings.Contains(file.Filename, ".proto") {
			schemaRequest.SchemaType = "PROTOBUF"
		}

		content, err := file.Open()
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(content)
		if err != nil {
			log.Fatal(err)
		}
		schemaRequest.SchemaBody = string(body)

		err = a.registerSchema(&schemaRequest)
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = "Failed to register the schema"
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		//Also send this to the content to the AvroExtractor to extract the fields
		fields := a.extractAvroFields(body)

		//Return the success response and the URL at which to publish the data to go to Kafka
		var response SchemaResponse
		response.Headers = []string{"--header 'Accept: application/vnd.kafka.v2+json, application/vnd.kafka+json, application/json'", "--header 'Content-Type: application/vnd.kafka.avro.v2+json'"}
		response.URL = fmt.Sprintf("http://restproxy/9082/topics/%s", schemaRequest.DatasetName)
		/*http://schemaregistry:8082/subjects/user-value/versions/latest/schema*/
		response.SchemaRegistryURL = fmt.Sprintf("http://schemaregistry:8082/schemas/%s-value/versions/latest/schema", schemaRequest.DatasetName)

		fmt.Println(fields)

		response.Fields = fields

		avroSample := AvroRequest{
			ValueSchemaID: schemaRequest.KafkaRegistryID,
			Records: []Record{
				Record{Value: []string{"Your avro formatted data compliant with registered schema"}},
			},
		}
		avroJSON, _ := json.Marshal(avroSample)
		response.SchemaBody = string(avroJSON)
		response.CurlCommand = fmt.Sprintf("curl --request POST --url %s  '%s'  '%s'  --data '%s'", response.URL, response.Headers[0], response.Headers[1], response.SchemaBody)

		fmt.Println(response.CurlCommand)

		c.JSON(http.StatusOK, response)
	}
}

type AvroExtractRequest struct {
	Schema string `json:"schema"`
}

func (a *App) extractAvroFields(content []byte) (successResponse AvroExtractResponse) {

	avroExtractRequest := AvroExtractRequest{Schema: string(content)}
	fmt.Println(avroExtractRequest)

	body, err := json.Marshal(avroExtractRequest)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:2602/demo/getFileContent", strings.NewReader(string(body)))
	if err != nil {
		// handle errr)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(respBody, &successResponse)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	log.Println("Succesfully registered the schema.")

	return
}

func (a *App) registerSchema(schemaRequest *SchemaRequest) (err error) {
	/*Once we have the file content, we need to do a few things.
	0. Get the datsetname so the corrsponding kafka topic can be created and registered
	1. Send this out to the schema registry to see if it's a valid request
	2. If and when it is registered, update the status in the metadata table to schema applied
	3. Insert record into the schemaregistry table with the corresponding ids
	4. Provide the end user with details on how to publish this data
	*/
	requestObj := KafkaSchemaRegistryRequest{Schema: schemaRequest.SchemaBody, SchemaType: schemaRequest.SchemaType}
	//Get the datasetname which will the Kafka topic name
	_, requestObj.DatasetName = a.getDatasetSetName(schemaRequest.CatalogDatasetID)
	if len(requestObj.DatasetName) == 0 {
		err = errors.New(fmt.Sprintf("Failed to find a dataset with the dataset id:", schemaRequest.CatalogDatasetID))
		return
	}
	schemaRequest.DatasetName = requestObj.DatasetName

	//Register schema to Kafka
	schemaRequest.KafkaRegistryID, err = a.registerSchemaToKafka(requestObj)
	if err != nil {
		return err
	}

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

	//Update the metadata tier status in the metadata table
	_, err = tx.Exec("UPDATE datacatalog.public.metadata set metadata_status=$1 where dataset_id=$2", "schema_applied", schemaRequest.CatalogDatasetID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		//To-do: De-register from schema registry
		return err
	}

	return
}

func (a *App) getDatasetSetName(id int64) (err error, datasetName string) {
	return a.DB.QueryRow("SELECT dataset_name FROM datacatalog.public.metadata where dataset_id=$1", id).Scan(&datasetName), datasetName
}

func (a *App) registerSchemaToKafka(request KafkaSchemaRegistryRequest) (id int64, err error) {
	id = int64(-1)
	payloadBytes, err := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://schemaregistry:8082/subjects/%s-value/versions", request.DatasetName), body)
	if err != nil {
		// handle errr)
		return
	}
	req.Header.Add("Accept", "application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var successResponse SchemaRegistrySuccessResponse
	var failureResponse SchemaRegistrySuccessResponse
	err = json.Unmarshal(respBody, &successResponse)
	if err != nil {
		err = json.Unmarshal(respBody, &failureResponse)
		if err != nil {
			err = errors.New("Invalid schema body. Unable to register schema")
			return
		}
		return
	}
	defer resp.Body.Close()

	log.Println("Succesfully registered the schema.")
	id = successResponse.Id

	return
}

//getSchema - Retrieves the schema from schema registry.
func (a *App) getSchema(schemaURL string) (schemaResponse SchemaRegistryResponse, err error) {
	//Let's get the schema from the schema registry

	fmt.Println("URL is ", schemaURL)

	req, err := http.NewRequest("GET", schemaURL, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		log.Fatal(err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { //
		log.Fatal(err)
		return
	}

	fmt.Println(respBody)

	err = json.Unmarshal(respBody, &schemaResponse)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	return
}

type SchemaRegistryResponse struct {
	Fields []struct {
		Doc  string        `json:"doc"`
		Name string        `json:"name"`
		Type []interface{} `json:"type"`
	} `json:"fields"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
}

type AvroExtractResponse struct {
	Status string `json:"status"`
	Fields []struct {
		Doc  string `json:"doc"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"data"`
}
