package main

import (
	"fmt"
	"strings"
	"time"
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
	IsComplexPrimaryKey          bool   `json:"is_complex_primary_key"`
}

type HiveSubmitPayload struct {
	PartitionedBy   string `json:"partitioned_by"`
	BasePath        string `json:"base_path"`
	Database        string `json:"database"`
	Table           string `json:"table"`
	IsMultiPartKeys bool   `json:"is_multi_part_keys"`
}

type HudiSyncSuccess struct {
	Message string `json:"message"`
}

//syncToHudi submits a spark-submit job to sync the topic contents to Hudi
func (a *App) syncToHudi(p PartitionDataset) (err error) {
	var successMessage HudiSyncSuccess

	sparkSubmitPayload := SparkSubmitPayload{
		Class:                        "org.apache.hudi.utilities.deltastreamer.HoodieDeltaStreamer",
		Jar:                          "/var/hoodie/ws/docker/hoodie/hadoop/hive_base/target/hoodie-utilities.jar",
		TableType:                    "MERGE_ON_READ",
		SourceClass:                  "org.apache.hudi.utilities.sources.AvroKafkaSource",
		SourceOrderingField:          "last_updated_time",
		TargetBasePath:               fmt.Sprintf("/user/hive/warehouse/%s_mor", p.DatasetName),
		TargetTable:                  fmt.Sprintf("%s_mor", p.DatasetName),
		SchemaProviderClass:          "org.apache.hudi.utilities.schema.SchemaRegistryProvider",
		HoodieConfSchemaRegistryURL:  fmt.Sprintf("hoodie.deltastreamer.schemaprovider.registry.url=http://schemaregistry:8082/subjects/%s-value/versions/latest", p.DatasetName),
		HoodieConfRecordKeyField:     fmt.Sprintf("hoodie.datasource.write.recordkey.field=%s", strings.Join(p.PrimaryKeys, ",")),
		HoodieConfPartitionPathField: fmt.Sprintf("hoodie.datasource.write.partitionpath.field=%s", p.PartitionPath),
		HoodieConfKafkaTopic:         fmt.Sprintf("hoodie.deltastreamer.source.kafka.topic=%s", p.DatasetName),
		HoodieBootStrapServers:       "bootstrap.servers=kafkabroker:9092",
		HoodieKafkaConsumerGroupID:   "group.id=hudi-deltastreamer-consumer",
		HoodieSchemaRegistyrURL:      "schema.registry.url=http://schemaregistry:8082",
		HoodieConfAutoOffsetReset:    "auto.offset.reset=earliest",
		HoodieStream:                 true,
	}

	hiveSubmitPayload := HiveSubmitPayload{
		PartitionedBy: p.PartitionPath,
		BasePath:      fmt.Sprintf("/user/hive/warehouse/%s_mor", p.DatasetName),
		Database:      "default",
		Table:         fmt.Sprintf("%s_mor", p.DatasetName),
	}

	if len(p.PrimaryKeys) > 1 {
		sparkSubmitPayload.IsComplexPrimaryKey = true
		hiveSubmitPayload.IsMultiPartKeys = true
	}

	headersMap := map[string]string{
		"Content-Type": "application/json",
	}
	sparkSubmitURL := "http://hudisync:2049/sparkSubmit"
	err = makePostCall(sparkSubmitPayload, headersMap, sparkSubmitURL, &successMessage, &successMessage)
	if err != nil {
		fmt.Println("Failed to submit the spark job")
		return
	}

	time.Sleep(30 * time.Second)

	hiveSubmitURL := "http://hudisync:2049/hiveSubmit"
	err = makePostCall(hiveSubmitPayload, headersMap, hiveSubmitURL, &successMessage, &successMessage)
	if err != nil {
		fmt.Println("Failed to sync to hive")
		return
	}

	fmt.Println("Returning from hudi sync function with the error ", err)

	return
}
