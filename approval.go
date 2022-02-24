package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) showReviewPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var metadataRequest MetadataRequest
		// Check if the dataset ID is valid
		if catalogDatasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64); err == nil {
			// Check if the dataset exists
			if dataset, err := a.getCompleteDatasetByID(catalogDatasetID); err == nil {

				a.getMetadata(catalogDatasetID, &metadataRequest)

				// Call the render function with the title, dataset and the name of the
				// template
				render(c, gin.H{
					"payload":  dataset,
					"metadata": metadataRequest,
				}, "approval.html")

			} else {
				// If the dataset is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}

		} else {
			// If an invalid dataset ID is specified in the URL, abort with an error
			c.AbortWithStatus(http.StatusNotFound)
		}

	}
}

func (a *App) getMetadata(datasetID int64, m *MetadataRequest) (err error) {
	err = a.DB.QueryRow("SELECT  dataset_name,  dataset_logical_name,  dataset_description,  dataset_type,  dataset_source,  dataset_share,  dataset_retention,  dataset_retention_justification,  dataset_arrival_frequency, organization,product,team,data_steward,platform_name, data_classification FROM datacatalog.public.metadata where dataset_id =$1", datasetID).Scan(&m.DatasetName, &m.DatasetLogicalName, &m.DatasetDescription, &m.DatasetType, &m.DatasetSource, &m.DatasetShare, &m.DatasetRetention, &m.DatasetRetentionJustification, &m.DatasetArrivalFrequency, &m.Organization, &m.Product, &m.Team, &m.DataSteward, &m.PlatformName, &m.DataClassiffication)
	return
}

func (a *App) approveDataset(datasetID int64) (err error) {
	_, err = a.DB.Exec("UPDATE datacatalog.public.metadata set datasteward_approved=true where dataset_id=$1", datasetID)
	return
}

func (a *App) approveDatasetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		resp, err := http.Get("http://127.0.0.1:9082/v3/clusters")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		var clusterResponse ClusterResponse
		err = json.Unmarshal([]byte(jsonDataFromHttp), &clusterResponse)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Cluster ID is ", clusterResponse.Data[0].ClusterID)
		clusterID := clusterResponse.Data[0].ClusterID

		if catalogDatasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64); err == nil {
			a.approveDataset(catalogDatasetID)
			datasetName := a.getDatasetName(catalogDatasetID)

			data := KafkaTopicPayload{
				TopicName:         datasetName,
				PartitionsCount:   1,
				ReplicationFactor: 1,
			}

			a.createKafkaTopic(data, clusterID)

		}

		//Get the Kafka Cluster Info
		c.JSON(http.StatusOK, gin.H{"status": "Approved"})
	}
}

func (a *App) getDatasetName(id int64) (datasetName string) {
	a.DB.QueryRow("Select dataset_name from datacatalog.public.metadata where dataset_id=$1", id).Scan(&datasetName)
	return
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
