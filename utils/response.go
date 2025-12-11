package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON mengirim respons JSON standar
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError mengirim respons error standar
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"status": "error", "message": message})
}