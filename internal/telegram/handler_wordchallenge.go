package telegram

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// handleWordChallenge generates a new Indonesian sentence challenge and sends it to the user
func (b *Bot) handleWordChallenge(chatID int64) {
	// Notify user that we're generating
	b.sendMessage(chatID, "⏳ Sedang menyiapkan soal\\.\\.\\.", "MarkdownV2")

	result, err := b.aiClient.GenerateSentence()
	if err != nil {
		b.sendMessage(chatID, "❌ Gagal membuat soal\\. Coba lagi dengan /word\\-challenge", "MarkdownV2")
		return
	}

	indonesian, _ := result["indonesian_sentence"].(string)
	correctAnswer, _ := result["correct_answer"].(string)
	topic, _ := result["topic"].(string)

	if indonesian == "" || correctAnswer == "" {
		b.sendMessage(chatID, "❌ Gagal membuat soal\\. Coba lagi dengan /word\\-challenge", "MarkdownV2")
		return
	}

	// Save challenge to session
	challengeID := uuid.New().String()
	b.sessions.Set(chatID, &PendingChallenge{
		ChallengeID:        challengeID,
		IndonesianSentence: indonesian,
		CorrectAnswer:      correctAnswer,
		Topic:              topic,
		CreatedAt:          time.Now(),
	})

	// Format and send the challenge
	topicLine := ""
	if topic != "" {
		topicLine = fmt.Sprintf("\n📌 *Topik:* %s", escapeMarkdown(topic))
	}

	msg := fmt.Sprintf(
		"🇮🇩 *Terjemahkan ke Business English:*\n\n_%s_%s\n\n💬 Ketik jawaban kamu\\, atau /skip untuk lewati\\.",
		escapeMarkdown(indonesian),
		topicLine,
	)
	b.sendMessage(chatID, msg, "MarkdownV2")
}

// handleAnswer evaluates the user's translation against the pending challenge
func (b *Bot) handleAnswer(chatID int64, userAnswer string) {
	challenge, ok := b.sessions.Get(chatID)
	if !ok {
		b.sendMessage(chatID, "Tidak ada soal aktif\\. Ketik /word\\-challenge untuk mulai\\.", "MarkdownV2")
		return
	}

	// Notify user that we're checking
	b.sendMessage(chatID, "🔍 Memeriksa jawaban\\.\\.\\.", "MarkdownV2")

	result, err := b.aiClient.CheckTranslation(challenge.IndonesianSentence, challenge.CorrectAnswer, userAnswer)
	if err != nil {
		b.sendMessage(chatID, "❌ Gagal memeriksa jawaban\\. Coba lagi\\.", "MarkdownV2")
		return
	}

	// Remove the session regardless of result
	b.sessions.Delete(chatID)

	isCorrect, _ := result["is_correct"].(bool)
	correctAnswer, _ := result["correct_answer"].(string)
	explanation, _ := result["explanation"].(string)
	corrections, _ := result["corrections"].(string)

	if correctAnswer == "" {
		correctAnswer = challenge.CorrectAnswer
	}

	// Build response message
	var statusLine string
	if isCorrect {
		statusLine = "✅ *Benar\\!* Terjemahan kamu tepat\\."
	} else {
		statusLine = "❌ *Kurang tepat\\.* Yuk perbaiki\\!"
	}

	msg := fmt.Sprintf("%s\n\n", statusLine)
	msg += fmt.Sprintf("🇮🇩 *Kalimat:*\n_%s_\n\n", escapeMarkdown(challenge.IndonesianSentence))
	msg += fmt.Sprintf("✍️ *Jawaban kamu:*\n`%s`\n\n", escapeMarkdown(userAnswer))
	msg += fmt.Sprintf("💡 *Jawaban terbaik:*\n`%s`\n", escapeMarkdown(correctAnswer))

	if explanation != "" {
		msg += fmt.Sprintf("\n📖 *Penjelasan:*\n%s\n", escapeMarkdown(explanation))
	}

	if corrections != "" && !isCorrect {
		msg += fmt.Sprintf("\n🔧 *Koreksi:*\n%s\n", escapeMarkdown(corrections))
	}

	msg += "\n▶️ Ketik /word\\-challenge untuk soal berikutnya\\."

	b.sendMessage(chatID, msg, "MarkdownV2")
}

// handleSkip cancels the current challenge
func (b *Bot) handleSkip(chatID int64) {
	challenge, ok := b.sessions.Get(chatID)
	if !ok {
		b.sendMessage(chatID, "Tidak ada soal aktif\\. Ketik /word\\-challenge untuk mulai\\.", "MarkdownV2")
		return
	}

	b.sessions.Delete(chatID)

	msg := fmt.Sprintf(
		"⏭️ Soal dilewati\\.\n\n💡 *Jawaban yang benar:*\n`%s`\n\nKetik /word\\-challenge untuk soal berikutnya\\.",
		escapeMarkdown(challenge.CorrectAnswer),
	)
	b.sendMessage(chatID, msg, "MarkdownV2")
}

// escapeMarkdown escapes special characters for Telegram MarkdownV2
func escapeMarkdown(s string) string {
	specialChars := []byte{'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'}
	result := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		for _, c := range specialChars {
			if s[i] == c {
				result = append(result, '\\')
				break
			}
		}
		result = append(result, s[i])
	}
	return string(result)
}
