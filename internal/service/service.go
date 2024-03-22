package service

//Единственный нужный импорт здесь :)
import "database/sql"

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

// GetCommands возвращзает список всех команд в бд
func (s *CommandService) GetCommands() ([]Command, error) {

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
	err := s.DB.QueryRow("SELECT id, content FROM commands WHERE id = $1", id).Scan(&command.ID, &command.Content)
	if err != nil {
		return nil, err
	}

	return &command, nil
}
