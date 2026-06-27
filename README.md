# Navisha English

Business English learning platform for IT professionals. Built with Next.js and Go.

## Overview

Navisha English helps developers and IT professionals improve their Business English from B1 to C1 level through five focused modules — all content is grounded in real IT workplace scenarios.

| Module | Description |
|---|---|
| Vocabulary Builder | 114+ IT business terms with spaced repetition (SM-2 algorithm) |
| Grammar | 13 exercises covering passive voice, conditionals, modal verbs, reported speech, and more |
| Business Writing | 17 exercises (email, proposal, report) with AI feedback on grammar and tone |
| Role-play Scenarios | 19 AI-powered conversations simulating standups, client calls, interviews, and code reviews |
| Word Challenge | AI-generated Indonesian→English translation challenges with instant feedback |

Additionally, a **Telegram Bot** is included for Word Challenge practice on the go.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 14 (App Router), TypeScript, Tailwind CSS, Zustand, TanStack Query |
| Backend | Go 1.26, Gin framework |
| Database | PostgreSQL |
| AI | DeepSeek API (grammar feedback, role-play, word challenges) |
| Auth | Google OAuth2 + JWT + refresh tokens |
| Bot | Telegram Webhook Bot |

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [PostgreSQL 14+](https://www.postgresql.org/download/)
- DeepSeek API key — get one at [platform.deepseek.com](https://platform.deepseek.com)
- Google OAuth credentials (see setup guide below)

---

## Getting Started

### 1. Clone and enter the project

```bash
git clone <your-repo-url>
cd navisha-english
```

### 2. Configure environment variables

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Edit `.env`:

```env
# Server
PORT=8010
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=navisha_english
DB_SSLMODE=disable

# JWT — use a long random string in production
JWT_SECRET=your-super-secret-jwt-key

# AI
DEEPSEEK_API_KEY=your_deepseek_api_key_here
DEEPSEEK_MODEL=deepseek-chat

# Frontend
FRONTEND_URL=http://localhost:3010

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id_here
GOOGLE_CLIENT_SECRET=your_google_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8010/api/v1/auth/google/callback

# Telegram (optional)
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com
```

### 3. Create the database

```bash
createdb navisha_english
```

Or using psql:

```sql
CREATE DATABASE navisha_english;
```

### 4. Run migrations and seed data

```bash
go run cmd/migrate/main.go
```

Expected output:

```
Applying 001_create_tables ... OK
Applying 002_seed_vocabulary ... OK
...
Applying 018_seed_scenarios_batch3 ... OK

18 migration(s) applied successfully.
```

### 5. Start the backend

```bash
go run cmd/server/main.go
```

Backend runs at `http://localhost:8010`

### 6. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs at `http://localhost:3010`

---

## Google OAuth Setup

### Step 1 — Create a Google Cloud Project

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Click the project dropdown → **New Project**
3. Name it `navisha-english` → **Create**

### Step 2 — Configure OAuth Consent Screen

1. Go to **APIs & Services** → **OAuth consent screen**
2. Select **External** → **Create**
3. Fill in App name, support email, and developer contact
4. Save and continue through all steps
5. Under **Test users**, add your own email

### Step 3 — Create OAuth Credentials

1. Go to **APIs & Services** → **Credentials** → **+ Create Credentials** → **OAuth Client ID**
2. Application type: **Web application**
3. Add **Authorized JavaScript origins**: `http://localhost:3010`
4. Add **Authorized redirect URIs**: `http://localhost:8010/api/v1/auth/google/callback`
5. Click **Create** and copy the **Client ID** and **Client Secret**
6. Paste them into `.env` as `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`

### OAuth Flow

```
User clicks "Continue with Google"
  → Frontend redirects to /api/v1/auth/google
    → Backend redirects to Google consent screen
      → User approves
        → Google redirects to /api/v1/auth/google/callback
          → Backend validates token, upserts user
            → Redirects to /auth/callback?access_token=...
              → Frontend stores tokens → redirects to /dashboard
```

---

## Telegram Bot Setup (Optional)

The Telegram bot allows users to practice Word Challenges directly from Telegram.

### Commands

| Command | Description |
|---|---|
| `/start` | Welcome message |
| `/word-challenge` | Start a new Indonesian→English challenge |
| `/skip` | Reveal the correct answer |
| `/help` | List available commands |

### Setup

1. Create a bot via [@BotFather](https://t.me/BotFather) and copy the token
2. Add `TELEGRAM_BOT_TOKEN` and `TELEGRAM_WEBHOOK_URL` to `.env`
3. The webhook registers automatically on server startup

---

## Project Structure

```
navisha-english/
├── cmd/
│   ├── server/
│   │   └── main.go              # HTTP server entrypoint
│   └── migrate/
│       └── main.go              # Migration runner CLI
│
├── internal/
│   ├── auth/
│   │   ├── handler.go           # Register, login, logout, refresh (email/password)
│   │   └── google.go            # Google OAuth2 handler
│   ├── ai/
│   │   └── client.go            # DeepSeek API client (5 AI functions)
│   ├── database/
│   │   └── database.go          # PostgreSQL connection
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication middleware
│   │   └── cors.go              # CORS middleware
│   ├── modules/
│   │   ├── vocabulary/
│   │   │   └── handler.go       # Vocabulary CRUD + SM-2 review
│   │   ├── grammar/
│   │   │   └── handler.go       # Grammar exercises + auto-grading
│   │   ├── writing/
│   │   │   └── handler.go       # Writing submissions + AI feedback
│   │   ├── roleplay/
│   │   │   └── handler.go       # Roleplay sessions + AI conversation
│   │   └── wordchallenge/
│   │       └── handler.go       # Indonesian→English challenge + AI check
│   └── telegram/
│       ├── bot.go               # Webhook handler + command dispatcher
│       ├── handler_wordchallenge.go  # Word challenge flow for Telegram
│       └── session.go           # In-memory session store (30 min TTL)
│
├── migrations/
│   ├── schema/
│   │   ├── 001_create_tables.sql          # All 10 core tables
│   │   ├── 010_add_google_oauth.sql       # No-op placeholder
│   │   ├── 020_create_word_challenge.sql  # Word challenge history table
│   │   └── 021_update_word_challenge.sql  # Update word challenge schema
│   └── seeds/
│       ├── 002_seed_vocabulary.sql
│       ├── 003_seed_grammar.sql
│       ├── 004_seed_writing.sql
│       ├── 005_seed_scenarios.sql
│       ├── 006_seed_vocabulary_additional.sql
│       ├── 007_seed_grammar_additional.sql
│       ├── 008_seed_writing_additional.sql
│       ├── 009_seed_scenarios_additional.sql
│       ├── 011_seed_vocabulary_batch2.sql
│       ├── 012_seed_grammar_batch2.sql
│       ├── 013_seed_writing_batch2.sql
│       ├── 014_seed_scenarios_batch2.sql
│       ├── 015_seed_vocabulary_batch3.sql
│       ├── 016_seed_grammar_batch3.sql
│       ├── 017_seed_writing_batch3.sql
│       └── 018_seed_scenarios_batch3.sql
│
├── frontend/
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/page.tsx       # Google Sign-In page
│   │   │   └── register/page.tsx    # Registration page
│   │   ├── auth/callback/page.tsx   # OAuth token handler
│   │   ├── dashboard/
│   │   │   ├── layout.tsx           # Protected sidebar layout
│   │   │   ├── page.tsx             # Dashboard with progress stats
│   │   │   ├── vocabulary/page.tsx  # Browse + SM-2 flashcard review
│   │   │   ├── grammar/page.tsx     # MCQ exercises with feedback
│   │   │   ├── writing/page.tsx     # Writing editor + AI feedback panel
│   │   │   ├── roleplay/page.tsx    # AI chat interface
│   │   │   └── word-challenge/page.tsx  # Translation challenge + history
│   │   ├── components/
│   │   │   └── ThemeToggle.tsx      # Dark/light mode toggle
│   │   ├── layout.tsx               # Root layout
│   │   ├── page.tsx                 # Public landing page
│   │   └── providers.tsx            # QueryClient + ThemeProvider + Toast
│   └── lib/
│       ├── api.ts                   # Axios client with JWT auto-refresh
│       ├── utils.ts                 # cn(), formatDate(), getDifficultyColor()
│       └── store/
│           └── auth.ts              # Zustand auth store with persistence
│
├── .env.example
├── .gitignore
└── go.mod
```

---

## Migration Runner

Migrations are tracked in the `navisha_english_migration` table. Each file runs exactly once and is recorded with a SHA-256 checksum.

```bash
# Apply all pending migrations (default)
go run cmd/migrate/main.go

# Explicit
go run cmd/migrate/main.go up

# Check migration status
go run cmd/migrate/main.go status
```

Status output example:

```
ID     Version    Name                          Status     Applied At
--------------------------------------------------------------------
1      001        create_tables                 applied    2026-06-27 06:00:00
2      002        seed_vocabulary               applied    2026-06-27 06:00:01
...
19     018        seed_scenarios_batch3         pending
```

To add a new migration, create a file with the naming pattern:

```
migrations/schema/022_add_new_table.sql
migrations/seeds/019_seed_new_data.sql
```

---

## API Reference

All protected routes require the `Authorization: Bearer <access_token>` header.

### Auth

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | Public | Create account with email/password |
| POST | `/api/v1/auth/login` | Public | Sign in with email/password |
| POST | `/api/v1/auth/refresh` | Public | Rotate access + refresh tokens |
| GET | `/api/v1/auth/google` | Public | Redirect to Google OAuth |
| GET | `/api/v1/auth/google/callback` | Public | Google OAuth callback |
| POST | `/api/v1/auth/google/verify` | Public | Verify Google ID token (One Tap) |
| POST | `/api/v1/auth/logout` | JWT | Sign out |
| GET | `/api/v1/auth/me` | JWT | Get current user profile |

### Vocabulary

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/vocabulary` | List all vocabulary (`?category=`) |
| GET | `/api/v1/vocabulary/review` | Get cards due for SM-2 review (max 20) |
| POST | `/api/v1/vocabulary/:id/review` | Submit review with quality score 0–5 |
| GET | `/api/v1/vocabulary/progress` | Learning statistics |

### Grammar

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/grammar/exercises` | List exercises (`?topic=&difficulty=`) |
| GET | `/api/v1/grammar/exercises/:id` | Get exercise with questions |
| POST | `/api/v1/grammar/exercises/:id/submit` | Submit answers, receive score + explanations |
| GET | `/api/v1/grammar/progress` | Completion statistics |

### Writing

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/writing/exercises` | List exercises (`?type=&difficulty=`) |
| GET | `/api/v1/writing/exercises/:id` | Get exercise with prompt and template |
| POST | `/api/v1/writing/exercises/:id/submit` | Submit writing, receive AI feedback |
| GET | `/api/v1/writing/progress` | Submission statistics |

### Role-play

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/roleplay/scenarios` | List all scenarios |
| GET | `/api/v1/roleplay/scenarios/:id` | Get scenario details |
| POST | `/api/v1/roleplay/scenarios/:id/start` | Start a new session |
| POST | `/api/v1/roleplay/sessions/:id/message` | Send a message |
| GET | `/api/v1/roleplay/sessions/:id` | Get session with full message history |

### Word Challenge

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/word-challenge/generate` | Generate a new Indonesian→English challenge |
| POST | `/api/v1/word-challenge/check` | Submit answer for AI evaluation |
| GET | `/api/v1/word-challenge/history` | Get last 20 challenge attempts |

### Telegram

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/telegram/webhook` | Telegram webhook receiver (public) |

---

## Database Tables

All tables use the `navisha_english_` prefix.

| Table | Description |
|---|---|
| `navisha_english_users` | User accounts (email, Google ID, level) |
| `navisha_english_refresh_tokens` | Active refresh token sessions |
| `navisha_english_vocabulary` | Vocabulary words with definitions |
| `navisha_english_user_vocabulary` | Per-user SM-2 spaced repetition state |
| `navisha_english_grammar_exercises` | Grammar exercises with JSONB question content |
| `navisha_english_user_grammar_progress` | Per-user grammar scores |
| `navisha_english_writing_exercises` | Writing exercises with prompts and templates |
| `navisha_english_user_writing_submissions` | User submissions with AI feedback JSONB |
| `navisha_english_roleplay_scenarios` | Roleplay scenario definitions |
| `navisha_english_roleplay_sessions` | Active/completed roleplay sessions with JSONB messages |
| `navisha_english_word_challenge_history` | Indonesian→English challenge results |
| `navisha_english_migration` | Migration tracking (version, checksum, applied_at) |

---

## Content Summary

| Module | Count | Levels |
|---|---|---|
| Vocabulary | 114 words | B1–C1 |
| Grammar exercises | 13 | B1–B2 |
| Writing exercises | 17 | B1–C1 |
| Role-play scenarios | 19 | B1–C1 |

### Vocabulary Categories
- Project Management (23 words)
- Technical Communication (22 words)
- Meeting & Email Phrases (23 words)
- DevOps & Cloud (15 words)

### Grammar Topics
- Passive Voice
- Modal Verbs
- Conditional Sentences
- Reported Speech
- Articles
- Present Perfect vs Simple Past
- Prepositions
- Tone & Register (Formal vs Informal)
- Gerunds & Infinitives
- Linking Words
- Phrasal Verbs
- Noun Phrases
- Future Forms

### Writing Types
- Email (bug reports, escalations, rejections, status updates, onboarding, Slack messages, LinkedIn posts)
- Proposal (technical proposals)
- Report (sprint reviews, incident reports, retrospectives, API documentation, self-assessments, meeting agendas)

### Role-play Scenarios
- Meeting & Standup (daily standup, sprint planning, QBR)
- Client Communication (deadline negotiation, outage explanation, vendor evaluation)
- Code & Technical (code review, presenting tech solutions, explaining technical debt)
- Career (job interview behavioral, job interview system design, salary negotiation, culture fit)
- Team Dynamics (onboarding, peer feedback, cross-team dependencies, remote work)
- Incident Response (production P0 war room)

---

## AI Features

Powered by **DeepSeek** (`deepseek-chat` by default). All AI calls go through the Go backend — the API key is never exposed to the frontend.

| Feature | What AI Does |
|---|---|
| Writing feedback | Scores writing 0–100, checks grammar, assesses tone, suggests vocabulary improvements, provides an improved version |
| Role-play conversation | Plays a realistic character (manager, client, interviewer) and responds contextually |
| Word challenge generation | Creates random Indonesian business IT sentences with English translations |
| Translation evaluation | Checks user's English translation for correctness, provides corrections and explanations |

---

## Environment Variables Reference

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Backend server port |
| `GIN_MODE` | No | `debug` | Set to `release` in production |
| `DB_HOST` | Yes | `localhost` | PostgreSQL host |
| `DB_PORT` | Yes | `5432` | PostgreSQL port |
| `DB_USER` | Yes | — | PostgreSQL username |
| `DB_PASSWORD` | Yes | — | PostgreSQL password |
| `DB_NAME` | Yes | — | PostgreSQL database name |
| `DB_SSLMODE` | No | `disable` | PostgreSQL SSL mode |
| `DATABASE_URL` | No | — | Full DSN (overrides individual DB_ vars) |
| `JWT_SECRET` | Yes | — | Secret key for signing JWT tokens |
| `DEEPSEEK_API_KEY` | Yes | — | DeepSeek API key |
| `DEEPSEEK_MODEL` | No | `deepseek-chat` | DeepSeek model name |
| `FRONTEND_URL` | No | `http://localhost:3000` | Allowed CORS origin |
| `GOOGLE_CLIENT_ID` | Yes* | — | Google OAuth client ID (*required for Google Sign-In) |
| `GOOGLE_CLIENT_SECRET` | Yes* | — | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | Yes* | — | OAuth callback URL |
| `TELEGRAM_BOT_TOKEN` | No | — | Telegram bot token (bot disabled if not set) |
| `TELEGRAM_WEBHOOK_URL` | No | — | Public URL for Telegram webhook registration |
