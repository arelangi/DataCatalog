package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *App) registerElasticSearchSinksHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errResp ErrorResponse

		datasetID, err := strconv.ParseInt(c.Param("dataset_id"), 10, 64)
		if err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return

		}

		//Now that we have the datset id, we can save this request to the db.
		if err = a.saveSinkRequestToDB(datasetID, "ELASTICSEARCH"); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Failed to save the request to database")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success"})

	}
}

func (a *App) saveSinkRequestToDB(datasetID int64, sink string) (err error) {
	_, err = a.DB.Exec("insert into datacatalog.public.sinks(dataset_id, sink) values ($1, $2)", datasetID, sink)
	return
}

func (a *App) getSinkByID(datasetID int64) (res Sinks, err error) {
	rows, err := a.DB.Query("select sink from datacatalog.public.sinks where dataset_id=$1", datasetID)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var sink string
		err = rows.Scan(&sink)
		if err != nil {
			return
		}
		res.SinkValues = append(res.SinkValues, Sink{SinkName: sink})
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

type Sinks struct {
	SinkValues []Sink
}

type Sink struct {
	SinkName string `json:"sink_name"`
}
