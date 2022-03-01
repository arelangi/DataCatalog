package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SparkSubmitPayload struct {
	Class                        string `json:"class"`
	Jar                          string `json:"jar"`
	TableType                    string `json:"table_type"`
	SourceClass                  string `json:"source_class"`
	SourceOrderingField          string `json:"source_ordering_field"`
	TargetBasePath               string `json:"target_base_path"`
	TargetTable                  string `json:"target_table"`
	SchemaProviderClass          string `json:"schema_provider_class"`
	HoodieConfSchemaRegistryURL  string `json:"hoodie_conf_schema_registry_url"`
	HoodieConfRecordKeyField     string `json:"hoodie_conf_record_key_field"`
	HoodieConfPartitionPathField string `json:"hoodie_conf_partition_path_field"`
	HoodieConfKafkaTopic         string `json:"hoodie_conf_kafka_topic"`
	HoodieBootStrapServers       string `json:"hoodie_boot_strap_servers"`
	HoodieKafkaConsumerGroupID   string `json:"hoodie_kafka_consumer_group_id"`
	HoodieSchemaRegistyrURL      string `json:"hoodie_schema_registyr_url"`
	HoodieConfAutoOffsetReset    string `json:"hoodie_conf_auto_offset_reset"`
	HoodieStream                 bool   `json:"hoodie_stream"`
}

type HudiSyncRequest struct {
	DatasetName         string
	Recordkey           string
	PartitionKey        string
	SourceOrderingField string
}

//syncToHudi submits a spark-submit job to sync the topic contents to Hudi
func (a *App) syncToHudi(p HudiSyncRequest) {
	fmt.Println("In sync to hudi")
	return
	payload := SparkSubmitPayload{
		Class:                        "org.apache.hudi.utilities.deltastreamer.HoodieDeltaStreamer",
		Jar:                          "/var/hoodie/ws/docker/hoodie/hadoop/hive_base/target/hoodie-utilities.jar",
		TableType:                    "MERGE_ON_READ",
		SourceClass:                  "org.apache.hudi.utilities.sources.AvroKafkaSource",
		SourceOrderingField:          p.SourceOrderingField,
		TargetBasePath:               fmt.Sprintf("/user/hive/warehouse/%s_mor", p.DatasetName),
		TargetTable:                  fmt.Sprintf("%s_mor", p.DatasetName),
		SchemaProviderClass:          "org.apache.hudi.utilities.schema.SchemaRegistryProvider",
		HoodieConfSchemaRegistryURL:  fmt.Sprintf("hoodie.deltastreamer.schemaprovider.registry.url=http://schemaregistry:8082/subjects/%s-value/versions/latest", p.DatasetName),
		HoodieConfRecordKeyField:     fmt.Sprintf("hoodie.datasource.write.recordkey.field=%s", p.Recordkey),
		HoodieConfPartitionPathField: fmt.Sprintf("hoodie.datasource.write.partitionpath.field=%s", p.PartitionKey),
		HoodieConfKafkaTopic:         fmt.Sprintf("hoodie.deltastreamer.source.kafka.topic=%s", p.DatasetName),
		HoodieBootStrapServers:       "bootstrap.servers=kafkabroker:9092",
		HoodieKafkaConsumerGroupID:   "group.id=hudi-deltastreamer-consumer",
		HoodieSchemaRegistyrURL:      "schema.registry.url=http://schemaregistry:8082",
		HoodieConfAutoOffsetReset:    "auto.offset.reset=earliest",
		HoodieStream:                 true,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://hudisync:9082/sparkSubmit", body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}
