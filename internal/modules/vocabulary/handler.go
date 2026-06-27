package vocabulary

import (
	"database/sql"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type Vocabulary struct {
	ID              string     `json:"id"`
	Word            string     `json:"word"`
	Definition      string     `json:"definition"`
	Category        string     `json:"category"`
	ExampleSentence string     `json:"example_sentence"`
	Difficulty      string     `json:"difficulty"`
	EaseFactor      float64    `json:"ease_factor,omitempty"`
	Interval        int        `json:"interval,omitempty"`
	Repetitions     int        `json:"repetitions,omitempty"`
	NextReview      *time.Time `json:"next_review,omitempty"`
}

type ReviewRequest struct {
	Quality int `json:"quality" binding:"required,min=0,max=5"`
}

// GetAll returns all vocabulary with user progress info
func (h *Handler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")
	category := c.Query("category")

	query := `
		SELECT v.id, v.word, v.definition, v.category, v.example_sentence, v.difficulty,
			COALESCE(uv.ease_factor, 2.5), COALESCE(uv.interval, 1),
			COALESCE(uv.repetitions, 0), uv.next_review
		FROM navisha_english_vocabulary v
		LEFT JOIN navisha_english_user_vocabulary uv ON v.id = uv.vocab_id AND uv.user_id = $1
		WHERE ($2 = '' OR v.category = $2)
		ORDER BY v.category, v.word`

	rows, err := h.db.Query(query, userID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vocabulary"})
		return
	}
	defer rows.Close()

	var items []Vocabulary
	for rows.Next() {
		var v Vocabulary
		if err := rows.Scan(&v.ID, &v.Word, &v.Definition, &v.Category,
			&v.ExampleSentence, &v.Difficulty, &v.EaseFactor, &v.Interval,
			&v.Repetitions, &v.NextReview); err != nil {
			continue
		}
		items = append(items, v)
	}

	if items == nil {
		items = []Vocabulary{}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read vocabulary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetDueForReview returns vocabulary due for spaced repetition review
func (h *Handler) GetDueForReview(c *gin.Context) {
	userID := c.GetString("user_id")

	// Get all vocab, prioritize: overdue first, then new (not yet started)
	query := `
		SELECT v.id, v.word, v.definition, v.category, v.example_sentence, v.difficulty,
			COALESCE(uv.ease_factor, 2.5), COALESCE(uv.interval, 1),
			COALESCE(uv.repetitions, 0), uv.next_review
		FROM navisha_english_vocabulary v
		LEFT JOIN navisha_english_user_vocabulary uv ON v.id = uv.vocab_id AND uv.user_id = $1
		WHERE uv.next_review IS NULL OR uv.next_review <= NOW()
		ORDER BY uv.next_review ASC NULLS FIRST
		LIMIT 20`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review items"})
		return
	}
	defer rows.Close()

	var items []Vocabulary
	for rows.Next() {
		var v Vocabulary
		if err := rows.Scan(&v.ID, &v.Word, &v.Definition, &v.Category,
			&v.ExampleSentence, &v.Difficulty, &v.EaseFactor, &v.Interval,
			&v.Repetitions, &v.NextReview); err != nil {
			continue
		}
		items = append(items, v)
	}

	if items == nil {
		items = []Vocabulary{}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read review items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items, "count": len(items)})
}

// SubmitReview processes SM-2 spaced repetition algorithm
func (h *Handler) SubmitReview(c *gin.Context) {
	userID := c.GetString("user_id")
	vocabID := c.Param("id")

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch current SM-2 state
	var easeFactor float64
	var interval, repetitions int
	err := h.db.QueryRow(
		`SELECT COALESCE(ease_factor, 2.5), COALESCE(interval, 1), COALESCE(repetitions, 0)
		FROM navisha_english_user_vocabulary WHERE user_id=$1 AND vocab_id=$2`,
		userID, vocabID,
	).Scan(&easeFactor, &interval, &repetitions)

	if err == sql.ErrNoRows {
		easeFactor = 2.5
		interval = 1
		repetitions = 0
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// SM-2 Algorithm
	q := float64(req.Quality)
	if q < 3 {
		// Failed — reset repetitions and interval, but still update ease factor
		repetitions = 0
		interval = 1
	} else {
		// Successful recall — advance interval
		if repetitions == 0 {
			interval = 1
		} else if repetitions == 1 {
			interval = 6
		} else {
			interval = int(math.Round(float64(interval) * easeFactor))
		}
		repetitions++
	}

	// Update ease factor for all reviews (SM-2 spec applies penalty on failure too)
	easeFactor = easeFactor + (0.1 - (5-q)*(0.08+(5-q)*0.02))
	if easeFactor < 1.3 {
		easeFactor = 1.3
	}

	nextReview := time.Now().Add(time.Duration(interval) * 24 * time.Hour)

	_, err = h.db.Exec(`
		INSERT INTO navisha_english_user_vocabulary (user_id, vocab_id, ease_factor, interval, repetitions, next_review, last_reviewed)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (user_id, vocab_id) DO UPDATE SET
			ease_factor=$3, interval=$4, repetitions=$5, next_review=$6, last_reviewed=NOW()`,
		userID, vocabID, easeFactor, interval, repetitions, nextReview,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ease_factor":  easeFactor,
		"interval":     interval,
		"repetitions":  repetitions,
		"next_review":  nextReview,
	})
}

// GetProgress returns vocabulary learning statistics
func (h *Handler) GetProgress(c *gin.Context) {
	userID := c.GetString("user_id")

	var total, learned, dueToday int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM navisha_english_vocabulary`).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM navisha_english_user_vocabulary WHERE user_id=$1 AND repetitions > 0`, userID).Scan(&learned); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM navisha_english_user_vocabulary WHERE user_id=$1 AND next_review <= NOW()`, userID).Scan(&dueToday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	progress := 0.0
	if total > 0 {
		progress = float64(learned) / float64(total) * 100
	}
	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"learned":   learned,
		"due_today": dueToday,
		"progress":  progress,
	})
}
