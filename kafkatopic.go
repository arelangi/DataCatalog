package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

//createKafkaTopic creates the kafka topic to publish messages to
func (a *App) createKafkaTopic(payload KafkaTopicPayload, clusterID string) (err error) {
	var succesMessage KafkaTopicSuccessMsg
	var failureMessage KafkaTopicFailureMsg

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal the payload to JSON")
		return
	}

	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://restproxy:9082/v3/clusters/%s/topics", clusterID), body)
	if err != nil {
		fmt.Println("Failed to create the http request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failure on executing the request to create kafka topic")
		return
	}
	defer resp.Body.Close()

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read the response from Kafka")
		return
	}

	err = json.Unmarshal([]byte(jsonDataFromHttp), &succesMessage)
	if err != nil {
		err = json.Unmarshal([]byte(jsonDataFromHttp), &failureMessage)
		if err != nil {
			fmt.Println("Unkown response returned from Kafka")
			return
		} else {
			if strings.Contains(failureMessage.Message, "already exists") {
				err = nil
			} else {
				fmt.Println("Unexpected response returned from Kafka")
				return
			}
		}
	}
	return
}

type KafkaTopicSuccessMsg struct {
	ClusterID string `json:"cluster_id"`
	Configs   struct {
		Related string `json:"related"`
	} `json:"configs"`
	IsInternal bool   `json:"is_internal"`
	Kind       string `json:"kind"`
	Metadata   struct {
		ResourceName string `json:"resource_name"`
		Self         string `json:"self"`
	} `json:"metadata"`
	PartitionReassignments struct {
		Related string `json:"related"`
	} `json:"partition_reassignments"`
	Partitions struct {
		Related string `json:"related"`
	} `json:"partitions"`
	PartitionsCount   int64  `json:"partitions_count"`
	ReplicationFactor int64  `json:"replication_factor"`
	TopicName         string `json:"topic_name"`
}

type KafkaTopicFailureMsg struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
