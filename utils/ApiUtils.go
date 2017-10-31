package utils

import (
	"encoding/json"
	"net/http"
)

func HandleResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\"is_error\": true, \"msg\": \"" + err.Error() + "\"}"));
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\"is_error\": true, \"msg\": \"" + err.Error() + "\"}"));
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func HandleStatus(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\"is_error\": true, \"msg\": \"" + err.Error() + "\"}"));
		return
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}