package config

import (
	"net"
	"os"
	"strings"
)

type Config struct {
	Port          string
	LogLevel      string
	SecurityToken string
	AllowedIPs    []net.IP
}

func Load() *Config {
	cfg := &Config{
		Port:          getEnvOrDefault("PORT", "8080"),
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "info"),
		SecurityToken: os.Getenv("SECURITY_TOKEN"),
	}

	// Parse allowed IPs if set
	if allowedIPs := os.Getenv("ALLOWED_IPS"); allowedIPs != "" {
		for _, ipStr := range strings.Split(allowedIPs, ",") {
			if ip := net.ParseIP(strings.TrimSpace(ipStr)); ip != nil {
				cfg.AllowedIPs = append(cfg.AllowedIPs, ip)
			}
		}
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
