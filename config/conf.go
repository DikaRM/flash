package config

import (
	"os"
)

type Config struct {
	GoogleClient string
	GoogleSecret string
}

func LoadConfig() *Config {
	return &Config{
		GoogleClient: "496850955192-hjjpra27hmnif3a7k8qnlou70r5s4rb1.apps.googleusercontent.com",
		GoogleSecret: os.Getenv("GOOGLE_SECRET_CLIENT"),
	}

}
