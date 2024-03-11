package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"myapp/internal/service"

	"github.com/gorilla/mux"
)

// CommandHandler represents the HTTP handler for commands.
type CommandHandler struct {
	CommandService *service.CommandService
}

// NewCommandHandler creates a new CommandHandler with the given service.
func NewCommandHandler(commandService *service.CommandService) *CommandHandler {
	return &CommandHandler{CommandService: commandService}
}

// CreateCommand creates a new command.
func (h *CommandHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	var command service.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.CommandService.CreateCommand(&command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(command)
}

// GetCommands returns the list of all commands.
func (h *CommandHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	commands, err := h.CommandService.GetCommands()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(commands)
}

// GetCommand returns a specific command by ID.
func (h *CommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid command ID", http.StatusBadRequest)
		return
	}

	command, err := h.CommandService.GetCommand(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Command not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(command)
}
