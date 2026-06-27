package wordchallenge

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"navisha-english/backend/internal/ai"
)

type Handler struct {
	db       *sql.DB
	aiClient *ai.Client
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:       db,
		aiClient: ai.NewClient(),
	}
}

// GenerateSentence calls AI to get a random Indonesian sentence for the user to translate
func (h *Handler) GenerateSentence(c *gin.Context) {
	sentence, err := h.aiClient.GenerateSentence()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate sentence: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sentence})
}

type CheckRequest struct {
	ChallengeID        string `json:"challenge_id" binding:"required"`
	IndonesianSentence string `json:"indonesian_sentence" binding:"required"`
	CorrectAnswer      string `json:"correct_answer" binding:"required"`
	UserAnswer         string `json:"user_answer" binding:"required"`
}

// CheckAnswer evaluates the user's translation and saves result to history
func (h *Handler) CheckAnswer(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	result, err := h.aiClient.CheckTranslation(req.IndonesianSentence, req.CorrectAnswer, req.UserAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check answer: " + err.Error()})
		return
	}

	// Extract fields safely
	isCorrect, _ := result["is_correct"].(bool)
	correctAnswer, _ := result["correct_answer"].(string)
	explanation, _ := result["explanation"].(string)
	corrections, _ := result["corrections"].(string)

	if correctAnswer == "" {
		correctAnswer = req.CorrectAnswer
	}

	// Save to history
	if _, err := h.db.Exec(`
		INSERT INTO navisha_english_word_challenge_history
			(user_id, challenge_id, indonesian_sentence, correct_answer, user_answer, is_correct, explanation, corrections)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, req.ChallengeID, req.IndonesianSentence, correctAnswer, req.UserAnswer, isCorrect, explanation, corrections,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_correct":     isCorrect,
		"correct_answer": correctAnswer,
		"user_answer":    req.UserAnswer,
		"explanation":    explanation,
		"corrections":    corrections,
	})
}

// ClearHistory deletes all word challenge history for the user
func (h *Handler) ClearHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	_, err := h.db.Exec(`
		DELETE FROM navisha_english_word_challenge_history
		WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History cleared"})
}

func (h *Handler) GetHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.db.Query(`
		SELECT indonesian_sentence, correct_answer, user_answer, is_correct, explanation, corrections, attempted_at
		FROM navisha_english_word_challenge_history
		WHERE user_id = $1
		ORDER BY attempted_at DESC
		LIMIT 20`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var indonesianSentence, correctAnswer, userAnswer, explanation, corrections, attemptedAt string
		var isCorrect bool
		if err := rows.Scan(&indonesianSentence, &correctAnswer, &userAnswer, &isCorrect, &explanation, &corrections, &attemptedAt); err != nil {
			continue
		}
		history = append(history, map[string]interface{}{
			"indonesian_sentence": indonesianSentence,
			"correct_answer":      correctAnswer,
			"user_answer":         userAnswer,
			"is_correct":          isCorrect,
			"explanation":         explanation,
			"corrections":         corrections,
			"attempted_at":        attemptedAt,
		})
	}

	if history == nil {
		history = []map[string]interface{}{}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
