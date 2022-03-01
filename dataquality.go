package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) dataQualityHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request DQRequest
		var errResp ErrorResponse

		if err := c.Bind(&request); err != nil {
			errResp.Error = err.Error()
			errResp.Message = fmt.Sprintf("Invalid request. Error in request body")
			c.JSON(http.StatusBadRequest, errResp)
			return
		}

		err := a.saveDataQualityRulestoDB(request)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.JSON(http.StatusOK, request)

	}
}

func (a *App) saveDataQualityRulestoDB(r DQRequest) (err error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return
	}
	for _, v := range r.Rules {
		_, err = tx.Exec("INSERT into datacatalog.public.dataquality(dataset_id, description, type, values, field_name) VALUES ($1, $2, $3, $4, $5) on conflict do nothing", r.DatasetID, v.Description, v.RuleType, v.Values, v.FieldName)
		if err != nil {
			return
		}
	}

	//Update the metadata tier status in the metadata table
	_, err = tx.Exec("UPDATE datacatalog.public.metadata set metadata_status=$1 where dataset_id=$2", "DQ Applied", r.DatasetID)
	if err != nil {
		return
	}

	return
}

type DQRequest struct {
	Rules     []DQRule `json:"rules"`
	DatasetID int64    `json:"dataset_id"`
}

type DQRule struct {
	Description string `json:"description"`
	FieldName   string `json:"field_name"`
	RuleType    string `json:"rule_type"`
	Values      string `json:"values"`
}

func (a *App) getDataQualityRulesByID(id int64) (dataset DQRequest, err error) {
	rows, err := a.DB.Query("select field_name, type, description, values from datacatalog.public.dataquality where dataset_id=$1", id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r DQRule
		err = rows.Scan(&r.FieldName, &r.RuleType, &r.Description, &r.Values)
		if err != nil {
			return
		}

		dataset.Rules = append(dataset.Rules, r)
	}

	err = rows.Err()
	if err != nil {
		return
	}
	dataset.DatasetID = id
	fmt.Println(dataset)

	return
}
