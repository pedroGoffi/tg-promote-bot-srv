package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const imageDir = "./images" // Diretório onde as imagens serão armazenadas

// HandleUpdates processa as mensagens recebidas do Telegram
func HandleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue // Ignora updates sem mensagem
		}

		user := update.Message.From.UserName
		text := update.Message.Text

		log.Printf("[TELEGRAM] Mensagem recebida de [%s]: %s", user, text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá! Como posso ajudar?")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("[ERRO] Falha ao enviar mensagem para [%s]: %v", user, err)
		}
	}
}

// DownloadImageFromTelegram baixa uma imagem de um link do Telegram
func DownloadImageFromTelegram(c *gin.Context) {
	var req struct {
		TelegramLink string `json:"telegram_link"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[ERRO] JSON inválido:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido"})
		return
	}

	log.Println("[INFO] Link recebido:", req.TelegramLink)

	imageName := extractImageName(req.TelegramLink)
	imagePath := filepath.Join(imageDir, imageName)

	if fileExists(imagePath) {
		log.Println("[INFO] Imagem já existe, retornando endpoint.")
		c.JSON(http.StatusOK, gin.H{"message": "Imagem já existente", "endpoint": "/static/" + imageName})
		return
	}

	if err := downloadImage(req.TelegramLink, imagePath); err != nil {
		log.Println("[ERRO] Falha ao baixar imagem:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao baixar imagem"})
		return
	}

	log.Println("[SUCESSO] Imagem baixada:", imagePath)
	c.JSON(http.StatusOK, gin.H{"message": "Imagem baixada com sucesso", "endpoint": "/static/" + imageName})
}

// extractImageName extrai o nome da imagem do link
func extractImageName(telegramLink string) string {
	parts := strings.Split(telegramLink, "/")
	return parts[len(parts)-1]
}

// fileExists verifica se um arquivo já existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FindImage percorre o HTML e retorna o valor do atributo src da imagem desejada
func FindImage(body *html.Node) (string, error) {
	var src string
	var search func(*html.Node)

	search = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "tgme_page_photo_image" {
					for _, a := range n.Attr {
						if a.Key == "src" {
							src = a.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil && src == ""; c = c.NextSibling {
			search(c)
		}
	}

	search(body)

	if src == "" {
		return "", fmt.Errorf("imagem não encontrada")
	}
	return src, nil
}

func parseResponse(resp *http.Response) (*html.Node, error) {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar HTML: %w", err)
	}

	return doc, nil
}

func downloadImagebase64(dataURI string) ([]byte, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid data URI format")
	}

	// Decodifica caracteres URL antes de processar Base64
	decodedData, err := url.QueryUnescape(parts[1])
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar URL: %w", err)
	}
	// Decode the Base64 part
	data, err := base64.StdEncoding.DecodeString(decodedData)
	if err != nil {
		return nil, fmt.Errorf("error decoding Base64: %w", err)
	}

	return data, nil
}

// downloadImage baixa uma imagem de um link fornecido
func downloadImage(url, filePath string) error {
	log.Println("[DOWNLOAD] Iniciando download:", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao link: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar imagem, status: %s", resp.Status)
	}

	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	doc, err := parseResponse(resp)
	if err != nil {
		return fmt.Errorf("erro ao realizar parser sobre o link %s", url)
	}

	link, err := FindImage(doc)
	if err != nil {
		return fmt.Errorf("falha ao realizar download de imagem sobre o link %s", url)
	}
	resp.Body.Close()
	var data io.Reader

	if strings.HasPrefix(link, "http") {
		resp, err := http.Get(link)
		if err != nil {
			return fmt.Errorf("falha ao capturar imagem: %w", err)
		}
		defer resp.Body.Close()

		data = resp.Body
	} else if strings.HasPrefix(link, "data:image") {
		decoded, err := downloadImagebase64(link)
		if err != nil {
			return fmt.Errorf("falha ao realizar decode de imagem base 64\n[%s]", link)
		}

		data = bytes.NewReader(decoded)
		filePath = filePath + ".svg"

	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, data); err != nil {
		return fmt.Errorf("erro ao salvar imagem: %w", err)
	}

	log.Println("[SUCESSO] Imagem salva:", filePath)
	return nil
}
