package service

import "database/sql"

// Command represents a command entity.
type Command struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// CommandService provides methods for managing commands.
type CommandService struct {
	DB *sql.DB
}

// NewCommandService creates a new CommandService with the given database connection.
func NewCommandService(db *sql.DB) *CommandService {
	return &CommandService{DB: db}
}

// CreateCommand creates a new command in the database.
func (s *CommandService) CreateCommand(command *Command) error {

	_, err := s.DB.Exec("INSERT INTO commands (content) VALUES ($1)", command.Content)
	return err
}

// GetCommands returns a list of all commands from the database.
func (s *CommandService) GetCommands() ([]Command, error) {

	rows, err := s.DB.Query("SELECT id, content FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		err := rows.Scan(&c.ID, &c.Content)
		if err != nil {
			return nil, err
		}
		commands = append(commands, c)
	}

	return commands, nil
}

// GetCommand returns a specific command by ID from the database.
func (s *CommandService) GetCommand(id int) (*Command, error) {
	var command Command
	err := s.DB.QueryRow("SELECT id, content FROM commands WHERE id = $1", id).Scan(&command.ID, &command.Content)
	if err != nil {
		return nil, err
	}

	return &command, nil
}
