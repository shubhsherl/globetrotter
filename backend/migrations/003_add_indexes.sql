-- Migration: 003_add_indexes.sql
-- Description: Add indexes for better performance

-- Add index on games.user_id
CREATE INDEX IF NOT EXISTS idx_games_user_id ON games(user_id);

-- Add index on game_questions.game_id
CREATE INDEX IF NOT EXISTS idx_game_questions_game_id ON game_questions(game_id);

-- Add index on users.username
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username); 