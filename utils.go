package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func makePostCall[Req any, SuccessResp any, FailureResp any](requestVariable Req, headers map[string]string, URL string, successResponseVariable *SuccessResp, failureResponseVariable *FailureResp) (err error) {
	body, err := json.Marshal(requestVariable)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", URL, strings.NewReader(string(body)))
	if err != nil {
		log.Println(err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(respBody, &successResponseVariable)
	if err != nil {
		err = json.Unmarshal(body, &failureResponseVariable)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}
