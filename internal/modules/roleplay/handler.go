package roleplay

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type Scenario struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Context        string   `json:"context"`
	AIRole         string   `json:"ai_role"`
	UserRole       string   `json:"user_role"`
	Difficulty     string   `json:"difficulty"`
	Tags           []string `json:"tags"`
}

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (h *Handler) GetScenarios(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, title, context, ai_role, user_role, difficulty, tags
		FROM navisha_english_roleplay_scenarios ORDER BY difficulty`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scenarios"})
		return
	}
	defer rows.Close()

	var scenarios []Scenario
	for rows.Next() {
		var s Scenario
		var tags []byte
		if err := rows.Scan(&s.ID, &s.Title, &s.Context, &s.AIRole, &s.UserRole, &s.Difficulty, &tags); err != nil {
			continue
		}
		if err := json.Unmarshal(tags, &s.Tags); err != nil {
			s.Tags = []string{}
		}
		scenarios = append(scenarios, s)
	}

	if scenarios == nil {
		scenarios = []Scenario{}
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read scenarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scenarios})
}

func (h *Handler) GetScenario(c *gin.Context) {
	id := c.Param("id")
	var s Scenario
	var tags []byte
	err := h.db.QueryRow(
		`SELECT id, title, context, ai_role, user_role, difficulty, tags FROM navisha_english_roleplay_scenarios WHERE id=$1`, id,
	).Scan(&s.ID, &s.Title, &s.Context, &s.AIRole, &s.UserRole, &s.Difficulty, &tags)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scenario not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if err := json.Unmarshal(tags, &s.Tags); err != nil {
		s.Tags = []string{}
	}
	c.JSON(http.StatusOK, s)
}

func (h *Handler) StartSession(c *gin.Context) {
	userID := c.GetString("user_id")
	scenarioID := c.Param("id")

	// Verify scenario exists
	var scenario Scenario
	var systemPrompt string
	err := h.db.QueryRow(
		`SELECT id, title, context, ai_role, user_role, ai_system_prompt FROM navisha_english_roleplay_scenarios WHERE id=$1`, scenarioID,
	).Scan(&scenario.ID, &scenario.Title, &scenario.Context, &scenario.AIRole, &scenario.UserRole, &systemPrompt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scenario not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get opening message from AI
	openingMsg, err := h.aiClient.RoleplayOpen(systemPrompt, scenario.Context, scenario.AIRole)
	if err != nil {
		openingMsg = "Hello! I'm ready to begin. How can I help you today?"
	}

	initialMessages := []Message{
		{Role: "assistant", Content: openingMsg, CreatedAt: time.Now()},
	}
	messagesJSON, err := json.Marshal(initialMessages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize messages"})
		return
	}

	sessionID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO navisha_english_roleplay_sessions (id, user_id, scenario_id, messages)
		VALUES ($1, $2, $3, $4)`,
		sessionID, userID, scenarioID, messagesJSON,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id": sessionID,
		"scenario":   scenario,
		"messages":   initialMessages,
	})
}

func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("id")

	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch session + scenario
	var scenarioID string
	var messagesRaw []byte
	err := h.db.QueryRow(
		`SELECT scenario_id, messages FROM navisha_english_roleplay_sessions WHERE id=$1 AND user_id=$2 AND status='active'`,
		sessionID, userID,
	).Scan(&scenarioID, &messagesRaw)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var systemPrompt, aiRole, context string
	err = h.db.QueryRow(
		`SELECT ai_system_prompt, ai_role, context FROM navisha_english_roleplay_scenarios WHERE id=$1`, scenarioID,
	).Scan(&systemPrompt, &aiRole, &context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scenario"})
		return
	}

	var messages []Message
	if err := json.Unmarshal(messagesRaw, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse session messages"})
		return
	}

	// Append user message
	userMsg := Message{Role: "user", Content: req.Content, CreatedAt: time.Now()}
	messages = append(messages, userMsg)

	// Build conversation history for AI
	history := make([]map[string]string, 0, len(messages))
	for _, m := range messages {
		history = append(history, map[string]string{"role": m.Role, "content": m.Content})
	}

	// Get AI response
	aiResponse, err := h.aiClient.RoleplayReply(systemPrompt, context, aiRole, history)
	if err != nil {
		aiResponse = "I apologize, I'm having trouble responding right now. Could you repeat that?"
	}

	aiMsg := Message{Role: "assistant", Content: aiResponse, CreatedAt: time.Now()}
	messages = append(messages, aiMsg)

	// Save updated messages
	updatedJSON, err := json.Marshal(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize messages"})
		return
	}
	if _, err := h.db.Exec(
		`UPDATE navisha_english_roleplay_sessions SET messages=$1, updated_at=NOW() WHERE id=$2`,
		updatedJSON, sessionID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  aiMsg,
		"messages": messages,
	})
}

func (h *Handler) GetSession(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("id")

	var scenarioID, status string
	var messagesRaw []byte
	var createdAt time.Time

	err := h.db.QueryRow(
		`SELECT scenario_id, messages, status, created_at FROM navisha_english_roleplay_sessions WHERE id=$1 AND user_id=$2`,
		sessionID, userID,
	).Scan(&scenarioID, &messagesRaw, &status, &createdAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var messages []Message
	if err := json.Unmarshal(messagesRaw, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse session messages"})
		return
	}
	if messages == nil {
		messages = []Message{}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          sessionID,
		"scenario_id": scenarioID,
		"messages":    messages,
		"status":      status,
		"created_at":  createdAt,
	})
}
