package web

import (
	"encoding/json"
	"net/http"
)

//JSONStandardResponse standard json response
func JSONStandardRespond(w http.ResponseWriter, body interface{}, status int) {
	//TODO log the message and status at detail level
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)

	data, err := json.Marshal(body)
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

type JSONErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

//JSONStandardResponse standard json response
func JSONErrorRespond(w http.ResponseWriter, message string, status int) {
	JSONStandardRespond(w, JSONErrorResponse{
		Status: status,
		Error:  message,
	}, status)
}
