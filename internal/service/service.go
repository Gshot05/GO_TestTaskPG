package service

//Единственный нужный импорт здесь :)
import (
	"database/sql"
	"fmt"
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

	_, err := s.DB.Exec("INSERT INTO commands (content) VALUES ($1)", command.Content)
	return err
}

// GetCommands возвращает список всех команд или сообщение, если список пуст
func (s *CommandService) GetCommands() (string, error) {
	rows, err := s.DB.Query("SELECT * FROM commands")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		err := rows.Scan(&c.ID, &c.Content)
		if err != nil {
			return "", err
		}
		commands = append(commands, c)
	}

	// Если список команд пуст, возвращаем специальное сообщение
	if len(commands) == 0 {
		return "В базе данных пока нет команд", nil
	}

	// Если список команд не пуст, возвращаем их в виде строки
	return fmt.Sprintf("%+v", commands), nil
}

// GetCommand возвращает определенную команду по ID
func (s *CommandService) GetCommand(id int) (*Command, error) {
	var command Command
	err := s.DB.QueryRow("SELECT id, content FROM commands WHERE id = $1", id).Scan(&command.ID, &command.Content)
	if err != nil {
		return nil, err
	}

	return &command, nil
}
