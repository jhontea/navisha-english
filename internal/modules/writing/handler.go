package writing

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"navisha-english/backend/internal/ai"
)

type Handler struct {
	db        *sql.DB
	aiClient  *ai.Client
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:       db,
		aiClient: ai.NewClient(),
	}
}

type Exercise struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Context     string   `json:"context"`
	Prompt      string   `json:"prompt"`
	Template    string   `json:"template,omitempty"`
	KeyPhrases  []string `json:"key_phrases,omitempty"`
	Difficulty  string   `json:"difficulty"`
}

type SubmitRequest struct {
	Content string `json:"content" binding:"required,min=10"`
}

func (h *Handler) GetExercises(c *gin.Context) {
	exerciseType := c.Query("type")
	difficulty := c.Query("difficulty")

	query := `
		SELECT id, title, type, context, prompt, difficulty
		FROM navisha_english_writing_exercises
		WHERE ($1 = '' OR type = $1)
		AND ($2 = '' OR difficulty = $2)
		ORDER BY difficulty, type`

	rows, err := h.db.Query(query, exerciseType, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}
	defer rows.Close()

	var exercises []map[string]interface{}
	for rows.Next() {
		var e struct {
			ID, Title, Type, Context, Prompt, Difficulty string
		}
		if err := rows.Scan(&e.ID, &e.Title, &e.Type, &e.Context, &e.Prompt, &e.Difficulty); err != nil {
			continue
		}
		exercises = append(exercises, map[string]interface{}{
			"id": e.ID, "title": e.Title, "type": e.Type,
			"context": e.Context, "prompt": e.Prompt, "difficulty": e.Difficulty,
		})
	}

	if exercises == nil {
		exercises = []map[string]interface{}{}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exercises})
}

func (h *Handler) GetExercise(c *gin.Context) {
	id := c.Param("id")
	var e Exercise
	var templateVal sql.NullString
	var keyPhrasesRaw []byte

	err := h.db.QueryRow(
		`SELECT id, title, type, context, prompt, template, key_phrases, difficulty
		FROM navisha_english_writing_exercises WHERE id=$1`, id,
	).Scan(&e.ID, &e.Title, &e.Type, &e.Context, &e.Prompt, &templateVal, &keyPhrasesRaw, &e.Difficulty)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if templateVal.Valid {
		e.Template = templateVal.String
	}
	if len(keyPhrasesRaw) > 0 {
		if err := json.Unmarshal(keyPhrasesRaw, &e.KeyPhrases); err != nil {
			e.KeyPhrases = []string{}
		}
	}
	c.JSON(http.StatusOK, e)
}

func (h *Handler) Submit(c *gin.Context) {
	userID := c.GetString("user_id")
	exerciseID := c.Param("id")

	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch exercise for context
	var exercise Exercise
	err := h.db.QueryRow(
		`SELECT id, title, type, context, prompt FROM navisha_english_writing_exercises WHERE id=$1`, exerciseID,
	).Scan(&exercise.ID, &exercise.Title, &exercise.Type, &exercise.Context, &exercise.Prompt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get AI feedback
	feedback, err := h.aiClient.CheckWriting(exercise.Type, exercise.Context, exercise.Prompt, req.Content)
	if err != nil {
		// Save without feedback if AI fails
		feedback = map[string]interface{}{
			"error": "AI feedback unavailable",
			"score": 0,
		}
	}

	score := 0
	if s, ok := feedback["score"].(float64); ok {
		score = int(s)
	}

	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize feedback"})
		return
	}

	// Save submission
	if _, err := h.db.Exec(`
		INSERT INTO navisha_english_user_writing_submissions (user_id, exercise_id, content, feedback, score)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, exerciseID, req.Content, feedbackJSON, score,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedback": feedback,
		"score":    score,
	})
}

func (h *Handler) GetProgress(c *gin.Context) {
	userID := c.GetString("user_id")

	var totalExercises, submitted int
	var avgScore sql.NullFloat64

	if err := h.db.QueryRow(`SELECT COUNT(*) FROM navisha_english_writing_exercises`).Scan(&totalExercises); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err := h.db.QueryRow(
		`SELECT COUNT(DISTINCT exercise_id), AVG(score) FROM navisha_english_user_writing_submissions WHERE user_id=$1`, userID,
	).Scan(&submitted, &avgScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_exercises": totalExercises,
		"submitted":       submitted,
		"avg_score":       avgScore.Float64,
	})
}
