-- Migration: Alter word challenge history table to support sentences instead of words

-- Drop old table if it exists with old schema, recreate with correct columns
DROP TABLE IF EXISTS navisha_english_word_challenge_history;

CREATE TABLE IF NOT EXISTS navisha_english_word_challenge_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES navisha_english_users(id) ON DELETE CASCADE,
    challenge_id VARCHAR(50) NOT NULL,
    indonesian_sentence TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    user_answer TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    explanation TEXT,
    corrections TEXT,
    attempted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS navisha_english_word_challenge_history_user_idx
    ON navisha_english_word_challenge_history(user_id);
