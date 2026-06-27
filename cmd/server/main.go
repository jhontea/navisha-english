package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"navisha-english/backend/internal/ai"
	"navisha-english/backend/internal/auth"
	"navisha-english/backend/internal/database"
	"navisha-english/backend/internal/middleware"
	"navisha-english/backend/internal/modules/grammar"
	"navisha-english/backend/internal/modules/roleplay"
	"navisha-english/backend/internal/modules/vocabulary"
	"navisha-english/backend/internal/modules/wordchallenge"
	"navisha-english/backend/internal/modules/writing"
	"navisha-english/backend/internal/telegram"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	// Auth routes (public)
	authHandler := auth.NewHandler(db)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.Refresh)

	// Google OAuth routes (public)
	googleHandler := auth.NewGoogleHandler(db)
	api.GET("/auth/google", googleHandler.Redirect)
	api.GET("/auth/google/callback", googleHandler.Callback)
	api.POST("/auth/google/verify", googleHandler.VerifyIDToken)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.Auth())

	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/auth/me", authHandler.Me)

	// Vocabulary module
	vocabHandler := vocabulary.NewHandler(db)
	protected.GET("/vocabulary", vocabHandler.GetAll)
	protected.GET("/vocabulary/review", vocabHandler.GetDueForReview)
	protected.POST("/vocabulary/:id/review", vocabHandler.SubmitReview)
	protected.GET("/vocabulary/progress", vocabHandler.GetProgress)

	// Grammar module
	grammarHandler := grammar.NewHandler(db)
	protected.GET("/grammar/exercises", grammarHandler.GetExercises)
	protected.GET("/grammar/exercises/:id", grammarHandler.GetExercise)
	protected.POST("/grammar/exercises/:id/submit", grammarHandler.Submit)
	protected.GET("/grammar/progress", grammarHandler.GetProgress)

	// Writing module
	writingHandler := writing.NewHandler(db)
	protected.GET("/writing/exercises", writingHandler.GetExercises)
	protected.GET("/writing/exercises/:id", writingHandler.GetExercise)
	protected.POST("/writing/exercises/:id/submit", writingHandler.Submit)
	protected.GET("/writing/progress", writingHandler.GetProgress)

	// Word Challenge module
	wordChallengeHandler := wordchallenge.NewHandler(db)
	protected.GET("/word-challenge/generate", wordChallengeHandler.GenerateSentence)
	protected.POST("/word-challenge/check", wordChallengeHandler.CheckAnswer)
	protected.GET("/word-challenge/history", wordChallengeHandler.GetHistory)

	// Roleplay module
	roleplayHandler := roleplay.NewHandler(db)
	protected.GET("/roleplay/scenarios", roleplayHandler.GetScenarios)
	protected.GET("/roleplay/scenarios/:id", roleplayHandler.GetScenario)
	protected.POST("/roleplay/scenarios/:id/start", roleplayHandler.StartSession)
	protected.POST("/roleplay/sessions/:id/message", roleplayHandler.SendMessage)
	protected.GET("/roleplay/sessions/:id", roleplayHandler.GetSession)

	// Telegram Bot (Word Challenge)
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgToken != "" {
		tgBot, err := telegram.NewBot(tgToken, ai.NewClient())
		if err != nil {
			log.Printf("Warning: failed to initialize Telegram bot: %v", err)
		} else {
			// Register webhook route (public — Telegram calls this endpoint)
			api.POST("/telegram/webhook", tgBot.HandleUpdate)

			// Register webhook with Telegram on startup
			webhookURL := os.Getenv("TELEGRAM_WEBHOOK_URL")
			if webhookURL != "" {
				if err := tgBot.RegisterWebhook(webhookURL); err != nil {
					log.Printf("Warning: failed to register Telegram webhook: %v", err)
				}
			} else {
				log.Println("Warning: TELEGRAM_WEBHOOK_URL not set, skipping webhook registration")
			}
		}
	} else {
		log.Println("TELEGRAM_BOT_TOKEN not set, Telegram bot disabled")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
