-- Migration: 002_add_migrations_table.sql
-- Description: Add migrations table to track applied migrations

-- Create migrations table if it doesn't exist
CREATE TABLE IF NOT EXISTS migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 