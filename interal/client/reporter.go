package client

import (
	"bytes"
	"encoding/json"
	"gin_demo/interal"
	logger "gin_demo/interal/log"
	"net/http"
)
var log = logger.Log

func ReportUserAction(topic string, content []byte) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"topic": topic,
		"content": content,
	})
	body := bytes.NewBuffer(postBody)
	resp, err := http.Post(interal.GetReporterURL(), "application/json", body)
	if err != nil {
		log.Errorf("An Error ReportUserAction %v", err.Error())
		return
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusOK {
		log.Info("reporter ok!")
	} else {
		log.Errorf("bad request<%d>: %v", resp.StatusCode, result["message"])
	}
}
