package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

//createKafkaTopic creates the kafka topic to publish messages to
func (a *App) createKafkaTopic(payload KafkaTopicPayload, clusterID string) (err error) {
	var succesMessage KafkaTopicSuccessMsg
	var failureMessage KafkaTopicFailureMsg

	headersMap := map[string]string{
		"Content-Type": "application/json",
	}

	url := fmt.Sprintf("http://restproxy:9082/v3/clusters/%s/topics", clusterID)

	fmt.Println("Tried to register topic at the url", url)

	err = makePostCall(payload, headersMap, url, &succesMessage, &failureMessage)

	fmt.Println("error when registering topic is ", err)

	return
}

func (a *App) getKafkaClusterID() string {
	var clusterResponse ClusterResponse
	resp, err := http.Get("http://127.0.0.1:9082/v3/clusters")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = json.Unmarshal([]byte(jsonDataFromHttp), &clusterResponse)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("fucket ID is ", clusterResponse.Data[0].ClusterID)

	return clusterResponse.Data[0].ClusterID
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

type PartitionDataset struct {
	DatasetName   string
	PrimaryKeys   []string
	PartitionPath string
}

type ClusterResponse struct {
	Data []struct {
		Acls struct {
			Related string `json:"related"`
		} `json:"acls"`
		BrokerConfigs struct {
			Related string `json:"related"`
		} `json:"broker_configs"`
		Brokers struct {
			Related string `json:"related"`
		} `json:"brokers"`
		ClusterID      string `json:"cluster_id"`
		ConsumerGroups struct {
			Related string `json:"related"`
		} `json:"consumer_groups"`
		Controller struct {
			Related string `json:"related"`
		} `json:"controller"`
		Kind     string `json:"kind"`
		Metadata struct {
			ResourceName string `json:"resource_name"`
			Self         string `json:"self"`
		} `json:"metadata"`
		PartitionReassignments struct {
			Related string `json:"related"`
		} `json:"partition_reassignments"`
		Topics struct {
			Related string `json:"related"`
		} `json:"topics"`
	} `json:"data"`
	Kind     string `json:"kind"`
	Metadata struct {
		Next interface{} `json:"next"`
		Self string      `json:"self"`
	} `json:"metadata"`
}
