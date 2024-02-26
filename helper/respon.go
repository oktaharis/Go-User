package helper

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload) 
	w.Header().Set("Content-Type", "application/json") // perubahan Add menjadi Set
	w.WriteHeader(code) // perubahan writeHeader menjadi WriteHeader
	w.Write(response)
}