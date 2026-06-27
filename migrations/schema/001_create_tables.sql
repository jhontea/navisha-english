-- Migration: Create all core tables with navisha_english_ prefix

CREATE TABLE IF NOT EXISTS navisha_english_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL DEFAULT '',
    google_id VARCHAR(255),
    level VARCHAR(10) DEFAULT 'B1',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS navisha_english_users_google_id_idx
    ON navisha_english_users(google_id) WHERE google_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS navisha_english_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_vocabulary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    word VARCHAR(255) NOT NULL,
    definition TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    example_sentence TEXT NOT NULL,
    difficulty VARCHAR(10) DEFAULT 'B1',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_user_vocabulary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    vocab_id UUID NOT NULL REFERENCES navisha_english_vocabulary(id) ON DELETE CASCADE,
    ease_factor FLOAT DEFAULT 2.5,
    interval INTEGER DEFAULT 1,
    repetitions INTEGER DEFAULT 0,
    next_review TIMESTAMPTZ DEFAULT NOW(),
    last_reviewed TIMESTAMPTZ,
    UNIQUE(user_id, vocab_id)
);

CREATE TABLE IF NOT EXISTS navisha_english_grammar_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    instruction TEXT NOT NULL,
    content JSONB NOT NULL,
    difficulty VARCHAR(10) DEFAULT 'B1',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_user_grammar_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES navisha_english_grammar_exercises(id) ON DELETE CASCADE,
    score INTEGER NOT NULL,
    completed_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, exercise_id)
);

CREATE TABLE IF NOT EXISTS navisha_english_writing_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    context TEXT NOT NULL,
    prompt TEXT NOT NULL,
    template TEXT,
    key_phrases TEXT[],
    difficulty VARCHAR(10) DEFAULT 'B1',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_user_writing_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES navisha_english_writing_exercises(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    feedback JSONB,
    score INTEGER,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_roleplay_scenarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    context TEXT NOT NULL,
    ai_system_prompt TEXT NOT NULL,
    ai_role VARCHAR(255) NOT NULL,
    user_role VARCHAR(255) NOT NULL,
    difficulty VARCHAR(10) DEFAULT 'B1',
    tags TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS navisha_english_roleplay_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    scenario_id UUID NOT NULL REFERENCES navisha_english_roleplay_scenarios(id) ON DELETE CASCADE,
    messages JSONB DEFAULT '[]',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
