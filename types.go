package main

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type AvroExtractRequest struct {
	Schema string `json:"schema"`
}

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
	URL               string              `json:"url"`
	Headers           []string            `json:"headers"`
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