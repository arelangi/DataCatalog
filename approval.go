package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) approveDataset(datasetID int64) (err error) {
	_, err = a.DB.Exec("UPDATE datacatalog.public.metadata set datasteward_approved=true where dataset_id=$1", datasetID)
	return
}

func (a *App) getDatasetName(id int64) (datasetName string) {
	a.DB.QueryRow("Select dataset_name from datacatalog.public.metadata where dataset_id=$1", id).Scan(&datasetName)
	return
}

/*
	approveDatasetHandler performs the following actions once a dataset is approved by the data steward

	1. Create a Kafka Topic
	2. Submit a Spark Submit job to GoSparkServer with the corresponding topic
	3. Submit a Hive sync job
	4. Submit to sinks
	5. Register everything to datahub
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

		//Add ElasticSearch Sink
		err = a.addElasticSearchSinkByName(partitionDetails.DatasetName)
		if err != nil {
			fmt.Println("Failed to add ElasticSearch Sink:", err)
			panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": err.Error()})
			return
		}

		//Register to Datahub
		err = a.RegisterDatasetToDatahub(catalogDatasetID)
		if err != nil {
			fmt.Println("Failed to register to Datahub:", err)
			panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Failure", "message": err.Error()})
			return
		}

		//Get the Kafka Cluster Info
		c.JSON(http.StatusOK, gin.H{"status": "Approved"})
	}
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
