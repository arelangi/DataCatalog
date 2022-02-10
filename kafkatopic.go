package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type KafkaTopicPayload struct {
	TopicName         string              `json:"topic_name"`
	PartitionsCount   int                 `json:"partitions_count"`
	ReplicationFactor int                 `json:"replication_factor"`
	Configs           []KafkaTopicConfigs `json:"configs"`
}

type KafkaTopicConfigs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (a *App) createKafkaTopic(payload KafkaTopicPayload, clusterID string) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:9082/v3/clusters/%s/topics", clusterID), body)
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
