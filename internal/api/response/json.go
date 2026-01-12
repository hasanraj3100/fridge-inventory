package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println("Error encoding response:", err)
	}
}

func ResponseWithError(w http.ResponseWriter, statusCode int, message string) {
	ResponseWithJSON(w, statusCode, map[string]string{"error": message})
}
