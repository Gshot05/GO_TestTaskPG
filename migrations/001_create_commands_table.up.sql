-- migrations/001_create_commands_table.sql
CREATE TABLE IF NOT EXISTS commands (
    id SERIAL PRIMARY KEY,
    content TEXT
);
