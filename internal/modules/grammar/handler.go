package grammar

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type Exercise struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Topic       string          `json:"topic"`
	Instruction string          `json:"instruction"`
	Content     json.RawMessage `json:"content"`
	Difficulty  string          `json:"difficulty"`
}

type SubmitRequest struct {
	Answers map[string]string `json:"answers" binding:"required"`
}

func (h *Handler) GetExercises(c *gin.Context) {
	topic := c.Query("topic")
	difficulty := c.Query("difficulty")

	query := `
		SELECT id, title, topic, instruction, difficulty
		FROM navisha_english_grammar_exercises
		WHERE ($1 = '' OR topic = $1)
		AND ($2 = '' OR difficulty = $2)
		ORDER BY difficulty, topic`

	rows, err := h.db.Query(query, topic, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}
	defer rows.Close()

	var exercises []map[string]interface{}
	for rows.Next() {
		var e struct {
			ID          string
			Title       string
			Topic       string
			Instruction string
			Difficulty  string
		}
		if err := rows.Scan(&e.ID, &e.Title, &e.Topic, &e.Instruction, &e.Difficulty); err != nil {
			continue
		}
		exercises = append(exercises, map[string]interface{}{
			"id":          e.ID,
			"title":       e.Title,
			"topic":       e.Topic,
			"instruction": e.Instruction,
			"difficulty":  e.Difficulty,
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
	err := h.db.QueryRow(
		`SELECT id, title, topic, instruction, content, difficulty FROM navisha_english_grammar_exercises WHERE id=$1`, id,
	).Scan(&e.ID, &e.Title, &e.Topic, &e.Instruction, &e.Content, &e.Difficulty)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
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

	// Fetch exercise content to check answers
	var contentRaw json.RawMessage
	err := h.db.QueryRow(
		`SELECT content FROM navisha_english_grammar_exercises WHERE id=$1`, exerciseID,
	).Scan(&contentRaw)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Parse content to get correct answers
	var content struct {
		Questions []struct {
			ID            string `json:"id"`
			CorrectAnswer string `json:"correct_answer"`
			Explanation   string `json:"explanation"`
		} `json:"questions"`
	}
	if err := json.Unmarshal(contentRaw, &content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse exercise"})
		return
	}

	// Grade answers
	correct := 0
	results := make([]map[string]interface{}, 0, len(content.Questions))
	for _, q := range content.Questions {
		userAnswer := req.Answers[q.ID]
		isCorrect := userAnswer == q.CorrectAnswer
		if isCorrect {
			correct++
		}
		results = append(results, map[string]interface{}{
			"id":             q.ID,
			"correct":        isCorrect,
			"your_answer":    userAnswer,
			"correct_answer": q.CorrectAnswer,
			"explanation":    q.Explanation,
		})
	}

	total := len(content.Questions)
	score := 0
	if total > 0 {
		score = (correct * 100) / total
	}

	// Save progress
	if _, err := h.db.Exec(`
		INSERT INTO navisha_english_user_grammar_progress (user_id, exercise_id, score)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, exercise_id) DO UPDATE SET score=$3, completed_at=NOW()`,
		userID, exerciseID, score,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score":   score,
		"correct": correct,
		"total":   total,
		"results": results,
	})
}

func (h *Handler) GetProgress(c *gin.Context) {
	userID := c.GetString("user_id")

	var totalExercises, completed int
	var avgScore sql.NullFloat64

	if err := h.db.QueryRow(`SELECT COUNT(*) FROM navisha_english_grammar_exercises`).Scan(&totalExercises); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err := h.db.QueryRow(`SELECT COUNT(*), AVG(score) FROM navisha_english_user_grammar_progress WHERE user_id=$1`, userID).
		Scan(&completed, &avgScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_exercises": totalExercises,
		"completed":       completed,
		"avg_score":       avgScore.Float64,
	})
}
