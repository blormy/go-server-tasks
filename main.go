// Простой HHTP-сервер для управления списком задач
// Позволяет получать список задач и добавлять новые задачи
//
// Как запустить через PowerShell:
// $body = '{"Task":"Пример задачи"}'
// Invoke-RestMethod -Uri "http://localhost:8080/tasks" -Method Post -Body $body -ContentType "application/json; charset=utf-8"

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Структура таска, только одно поле - название
type Task struct {
	Title string `json:"Task"`
}

// Глобальный срез со всеми задачами
var tasks []Task

// Обрабатывает получение задач
// Возвращает список задач в формате JSON
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tasks)
}

// Обрабатывает добавление задач
// Ждет JSON с полем Task. Если поле пустое, то возвращает ошибку
func createTask(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newTask Task
	json.Unmarshal(body, &newTask)
	if newTask.Title == "" {
		http.Error(w, "Поле Task не может быть пустым", http.StatusBadRequest)
		return
	}
	tasks = append(tasks, newTask)
	fmt.Fprintln(w, "Задача добавлена!")
}

// Регестрирует обработчик и запускает сервер
func main() {

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTasks(w, r)
		} else if r.Method == "POST" {
			createTask(w, r)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
