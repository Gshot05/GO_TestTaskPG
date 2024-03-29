package service

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Command представляет собой структуру таблицы commands в базе данных
type Command struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// CommandService предоставляет собой методы для управления командами
type CommandService struct {
	DB *sql.DB
}

// NewCommandService создает новый CommandService с подключением к базе данных
func NewCommandService(db *sql.DB) *CommandService {
	return &CommandService{DB: db}
}

// CreateCommand создаёт новую команду в базе данных
func (s *CommandService) CreateCommand(command *Command) error {

	// SQL запрос на выборку данных
	_, err := s.DB.Exec("INSERT INTO commands (content) VALUES ($1)", command.Content)
	return err
}

// GetCommands возвращает список всех команд или сообщение, если список пуст
func (s *CommandService) GetCommands() ([]Command, error) {

	// SQL запрос на выборку данных
	rows, err := s.DB.Query("SELECT * FROM commands")
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

// GetCommand возвращает определенную команду по ID
func (s *CommandService) GetCommand(id int) (*Command, error) {
	var command Command
	// SQL запрос на выборку данных
	err := s.DB.QueryRow("SELECT id, content FROM commands WHERE id = $1", id).Scan(&command.ID, &command.Content)
	if err != nil {
		return nil, err
	}

	return &command, nil
}
