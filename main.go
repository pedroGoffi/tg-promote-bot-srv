package main

import (
	"bot-manager/internal/bot"
	"bot-manager/internal/server"
	"flag"
	"fmt"
)

func main() {
	cliMode := flag.Bool("cli", false, "Rodar em modo CLI")
	flag.Parse()

	bot.IniciarBot()

	if *cliMode {
		fmt.Println("Modo CLI ativado")
	} else {
		fmt.Println("Iniciando servidor...")
		server.StartServer()
	}
}
