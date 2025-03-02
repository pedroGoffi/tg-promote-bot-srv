package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig carrega as variáveis de ambiente do .env
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
}

// GetBotToken retorna o token do bot armazenado no .env
func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}
