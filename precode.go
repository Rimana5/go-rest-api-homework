package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func GetTasks(res http.ResponseWriter, req *http.Request) {
	tasks, err := json.Marshal(tasks)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json, charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(tasks)
	if err != nil {
		log.Println(err)
	}
}

func PostTasks(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := tasks[task.ID]; exists {
		errMsg := fmt.Sprintf("Task with ID %s already exists", task.ID)
		log.Println(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json, charset=utf-8")
	res.WriteHeader(http.StatusCreated)
}

func GetTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		log.Println("Task not found")
		http.Error(res, "Task not found", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json, charset=utf-8")
	res.WriteHeader(http.StatusCreated)
	if _, err = res.Write(resp); err != nil {
		log.Println(err)
	}
}

func DeleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		log.Println("Task not found")
		http.Error(res, "Task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	res.Header().Set("Content-Type", "application/json, charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", GetTasks)
	r.Post("/tasks", PostTasks)
	r.Post("/tasks/{id}", GetTask)
	r.Delete("/tasks/{id}", DeleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
