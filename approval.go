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

/*
	approveDatasetHandler performs the following actions once a dataset is approved by the data steward

	1. Create a Kafka Topic
	2. Submit a Spark Submit job to GoSparkServer with the corresponding topic
	3. Submit a Hive sync job
	4. Register the downstream dataset to the catalog
*/
func (a *App) approveDatasetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID := a.getKafkaClusterID()
		catalogDatasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": "Unable to parse the provided dataset id"})
			return
		}
		//Approve the dataset and get it's name
		err = a.approveDataset(catalogDatasetID)
		if err != nil {
			panic(err)
		}
		partitionDetails, err := a.getPartitionDetailsForDataset(catalogDatasetID)
		partitionDetails.DatasetName = a.getDatasetName(catalogDatasetID)

		if err != nil {
			panic(err)
		}
		//Create Kafka Topic
		KafkaTopicPayloadData := KafkaTopicPayload{
			TopicName:         partitionDetails.DatasetName,
			PartitionsCount:   1,
			ReplicationFactor: 1,
		}

		err = a.createKafkaTopic(KafkaTopicPayloadData, clusterID)
		if err != nil {
			fmt.Println("Failed to create kafka topic")
			panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": err.Error()})
			return
		}

		//Submit job to spark
		err = a.syncToHudi(partitionDetails)
		if err != nil {
			fmt.Println("Failed to submit hive job")
			panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": err.Error()})
			return
		}

		//Get the Kafka Cluster Info
		c.JSON(http.StatusOK, gin.H{"status": "Approved"})
	}
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

func (a *App) getDatasetName(id int64) (datasetName string) {
	a.DB.QueryRow("Select dataset_name from datacatalog.public.metadata where dataset_id=$1", id).Scan(&datasetName)
	return
}

func (a *App) getPartitionDetailsForDataset(id int64) (resp PartitionDataset, err error) {
	rows, err := a.DB.Query("select name from datacatalog.public.fields where dataset_id=$1 and primarykeyfield=true", id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}

		resp.PrimaryKeys = append(resp.PrimaryKeys, name)
	}

	err = a.DB.QueryRow("select name from datacatalog.public.fields where dataset_id=$1 and partitionfield=true", id).Scan(&resp.PartitionPath)
	return

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
