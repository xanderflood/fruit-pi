package web

import (
	"encoding/json"
	"net/http"
)

type JSONStandardResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//JSONStandardResponse standard json response
func JSONStandardRespond(w http.ResponseWriter, message string, status int) {
	//TODO log the message and status at detail level
	w.WriteHeader(status)

	data, err := json.Marshal(JSONStandardResponse{status, message})
	if err != nil {
		//TODO log
		return
	}

	_, err = w.Write(data)
	if err != nil {
		//TODO log
		return
	}
}
