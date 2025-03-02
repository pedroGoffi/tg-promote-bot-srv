package config

import (
	"os"
)

// GetBotToken retorna o token do bot armazenado no .env
func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}
