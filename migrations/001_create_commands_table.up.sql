-- migrations/001_create_commands_table.sql
CREATE TABLE IF NOT EXISTS commands (
    id SERIAL NOT NULL constraint PK_commands PRIMARY KEY,
    content TEXT
);