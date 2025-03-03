package config

import (
	"os"
)

// GetBotToken retorna o token do bot armazenado no .env
func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}

func GetServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "4000" // Default port
	}
	return port
}

func GetUT_KEY() string {
	return os.Getenv("UPLOADTHING_TOKEN")
}
