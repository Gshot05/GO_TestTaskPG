package handler

//Все необходимые библиотеки, а также импортируем внутренние зависимости проекта
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"myapp/internal/service"

	"github.com/gorilla/mux"
)

// CommandHandler представляет обработчик HTTP для команд
type CommandHandler struct {
	CommandService *service.CommandService
}

// NewCommandHandler создает новый CommandHandler с данным сервисом
func NewCommandHandler(commandService *service.CommandService) *CommandHandler {
	return &CommandHandler{CommandService: commandService}
}

// CreateCommand создаёт новую команду
func (h *CommandHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	var command service.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Обращаемся к команде из сервиса
	err := h.CommandService.CreateCommand(&command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(command)
}

// GetCommands возвращает список всех команд в бд
func (h *CommandHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	commands, err := h.CommandService.GetCommands()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если список команд пуст (возвращается ошибка с сообщением о пустом списке), возвращаем ошибку
	if commands == nil {
		http.Error(w, "В базе данных пока нет команд", http.StatusNotFound)
		return
	}

	// Если список команд не пуст, возвращаем их в виде JSON
	json.NewEncoder(w).Encode(commands)
}

// GetCommand возвращает определенную команду по ID
func (h *CommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "некорректный ID", http.StatusBadRequest)
		return
	}

	// Обращаемся к команде из сервиса
	command, err := h.CommandService.GetCommand(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Команда не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(command)
}
