package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yourorg/iam-credentials-proxy/internal/credentials"
)

func CredentialsHandler(w http.ResponseWriter, r *http.Request) {
	creds, err := credentials.Get()
	if err != nil {
		http.Error(w, "Failed to get credentials", http.StatusInternalServerError)
		log.Printf("Credential error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(creds)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
