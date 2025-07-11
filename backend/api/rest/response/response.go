package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//resp := APIResponse{
	//	Status: status,
	//	Data:   data,
	//}

	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, errMsg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Error: errMsg,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
