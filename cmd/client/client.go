package main

//Все необходимые библиотеки
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Создаём структуру Command которая повторяет таблицу commands в базе данных
type Command struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func main() {
	// Далее это опции для нашего клиента
	for {
		//Варианты событий
		fmt.Println("Добро пожаловать в терминал! Выберете команду:")
		fmt.Println("1. Получить все команды")
		fmt.Println("2. Получить команду по ID")
		fmt.Println("3. Добавить команду")
		fmt.Println("Для выхода введите 'exit'")

		// Тут мы выводим > и ждём какой будет выбор
		var choice string
		fmt.Print("> ")
		fmt.Scanln(&choice)

		// Switch-case для наших выборов
		switch choice {
		case "1":
			getAllCommands()
		case "2":
			var id int
			fmt.Print("Введите ID команды: ")
			fmt.Scanln(&id)
			getCommandByID(id)
		case "3":
			var content string
			fmt.Print("Введите содержание команды: ")
			fmt.Scanln(&content)
			addCommand(content)
		case "exit":
			fmt.Println("Программа завершена.")
			os.Exit(0)
		default:
			fmt.Println("Неправильный ввод")
		}
	}
}

// Отправляет HTTP запрос для нашего апи по нужному роуту, на получение всех команд
func getAllCommands() {
	resp, err := http.Get("http://localhost:8080/api/v1/commands")
	if err != nil {
		fmt.Println("Ошибка при выполнении GET запроса:", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Println("Ответ сервера:\n", buf.String())
}

// Отправляет HTTP запрос для нашего апи по нужному роуту, для получения команды по определённому ID
func getCommandByID(id int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/v1/commands/%d", id))
	if err != nil {
		fmt.Println("Ошибка при выполнении GET запроса:", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	fmt.Println("Ответ сервера:\n", buf.String())
}

// Отправляет HTTP запрос для нашего апи по нужному роуту, для добавления команды
func addCommand(content string) {
	command := Command{Content: content}
	jsonData, err := json.Marshal(command)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/v1/commands", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при выполнении POST запроса:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Команда успешно добавлена!")
}
