package bot

import (
	"bot-manager/internal/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IniciarBot() {
	// Obter token do bot
	botToken := config.GetBotToken()

	// Inicializar bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("🤖 Bot %s iniciado com sucesso!", bot.Self.UserName)
}
