package main

import "fmt"

type ElasticSearchPayload struct {
	Name   string              `json:"name"`
	Config ElasticSearchConfig `json:"config"`
}

type ElasticSearchConfig struct {
	ConnectorClass                  string `json:"connector.class"`
	ConnectionURL                   string `json:"connection.url"`
	ValueConverter                  string `json:"value.converter"`
	ValueConverterSchemaRegistryURL string `json:"value.converter.schema.registry.url"`
	TypeName                        string `json:"type.name"`
	Topics                          string `json:"topics"`
	KeyIgnore                       string `json:"key.ignore"`
	SchemaIgnore                    string `json:"schema.ignore"`
}

type ElasticSearchSuccessMessage struct {
	Config struct {
		Connection_url                      string `json:"connection.url"`
		Connector_class                     string `json:"connector.class"`
		Key_ignore                          string `json:"key.ignore"`
		Name                                string `json:"name"`
		Schema_ignore                       string `json:"schema.ignore"`
		Topics                              string `json:"topics"`
		Type_name                           string `json:"type.name"`
		Value_converter                     string `json:"value.converter"`
		Value_converter_schema_registry_url string `json:"value.converter.schema.registry.url"`
	} `json:"config"`
	Name  string   `json:"name"`
	Tasks []string `json:"tasks"`
	Type  string   `json:"type"`
}

type ElasticSearchFailureMessage struct {
	ErrorCode int64  `json:"error_code"`
	Message   string `json:"message"`
}

func (a *App) addElasticSearchSinkByName(datasetName string) (err error) {
	var successMessage ElasticSearchSuccessMessage
	var failureMessage ElasticSearchFailureMessage

	elasticSearchPayload := ElasticSearchPayload{
		Name: fmt.Sprintf("SINK_ELASTIC_AUTOMATIC_%s"),
		Config: ElasticSearchConfig{
			ConnectorClass:                  "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
			ConnectionURL:                   "http://elasticsearch:9200",
			ValueConverter:                  "io.confluent.connect.avro.AvroConverter",
			ValueConverterSchemaRegistryURL: "http://schemaregistry:8082",
			TypeName:                        "_doc",
			Topics:                          datasetName,
			KeyIgnore:                       "true",
			SchemaIgnore:                    "true",
		},
	}

	url := "http://elasticsearch:8083/connectors"
	headersMap := map[string]string{
		"Content-Type": "application/json",
	}
	err = makePostCall(elasticSearchPayload, headersMap, url, &successMessage, &failureMessage)
	if err != nil {
		fmt.Println("Failed to add elasticsearch connector:", err)
		return
	}

	return
}
