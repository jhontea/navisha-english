package telegram

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gin-gonic/gin"
	"navisha-english/backend/internal/ai"
)

// Bot holds the Telegram bot instance, AI client, and session store
type Bot struct {
	api      *tgbotapi.BotAPI
	aiClient *ai.Client
	sessions *SessionStore
}

// NewBot creates a new Telegram bot. Returns an error if the token is invalid.
func NewBot(token string, aiClient *ai.Client) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	log.Printf("Telegram bot authorized as @%s", botAPI.Self.UserName)

	return &Bot{
		api:      botAPI,
		aiClient: aiClient,
		sessions: NewSessionStore(),
	}, nil
}

// RegisterWebhook tells Telegram to send updates to the given URL.
// webhookURL should be the full public HTTPS URL, e.g.:
//
//	https://yourdomain.com/api/v1/telegram/webhook
func (b *Bot) RegisterWebhook(webhookURL string) error {
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("failed to build webhook config: %w", err)
	}

	resp, err := b.api.Request(wh)
	if err != nil {
		return fmt.Errorf("failed to register webhook: %w", err)
	}
	if !resp.Ok {
		return fmt.Errorf("telegram rejected webhook: %s", resp.Description)
	}

	info, err := b.api.GetWebhookInfo()
	if err == nil {
		log.Printf("Telegram webhook registered: %s", info.URL)
	}
	return nil
}

// HandleUpdate is the Gin handler for POST /api/v1/telegram/webhook
func (b *Bot) HandleUpdate(c *gin.Context) {
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Only handle regular text messages
	if update.Message == nil || update.Message.Text == "" {
		c.Status(http.StatusOK)
		return
	}

	msg := update.Message
	chatID := msg.Chat.ID
	text := msg.Text

	switch {
	case text == "/start":
		b.sendMessage(chatID, "Halo\\! 👋 Saya bot latihan Business English\\.\n\nKetik /word\\-challenge untuk mulai latihan terjemahan kalimat Indonesia ke English\\.", tgbotapi.ModeMarkdownV2)

	case text == "/word-challenge" || text == "/word_challenge":
		b.handleWordChallenge(chatID)

	case text == "/skip":
		b.handleSkip(chatID)

	case text == "/help":
		b.sendMessage(chatID, "Perintah yang tersedia:\n\n/word\\-challenge \\- Mulai soal terjemahan baru\n/skip \\- Lewati soal saat ini\n/help \\- Tampilkan bantuan ini", tgbotapi.ModeMarkdownV2)

	default:
		// Check if there's a pending challenge waiting for an answer
		if _, hasPending := b.sessions.Get(chatID); hasPending {
			b.handleAnswer(chatID, text)
		} else {
			b.sendMessage(chatID, "Ketik /word\\-challenge untuk mulai latihan\\.", tgbotapi.ModeMarkdownV2)
		}
	}

	c.Status(http.StatusOK)
}

// sendMessage sends a plain or markdown text message to the given chat
func (b *Bot) sendMessage(chatID int64, text string, parseMode ...string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if len(parseMode) > 0 {
		msg.ParseMode = parseMode[0]
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("telegram: failed to send message to chat %d: %v", chatID, err)
	}
}
