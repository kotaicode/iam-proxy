package credentials

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

type CredentialResponse struct {
	Version         int    `json:"Version"`
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

var cachedCreds *CredentialResponse
var cachedUntil time.Time

func Get() (*CredentialResponse, error) {
	// refresh if not cached or expired
	if cachedCreds == nil || time.Until(cachedUntil) < 5*time.Minute {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, err
		}
		creds, err := cfg.Credentials.Retrieve(context.Background())
		if err != nil {
			return nil, err
		}
		expiry := time.Now().Add(1 * time.Hour) // AWS returns actual expiration only for STS clients
		cachedCreds = &CredentialResponse{
			Version:         1,
			AccessKeyId:     creds.AccessKeyID,
			SecretAccessKey: creds.SecretAccessKey,
			SessionToken:    creds.SessionToken,
			Expiration:      expiry.UTC().Format(time.RFC3339),
		}
		cachedUntil = expiry
	}
	return cachedCreds, nil
}
