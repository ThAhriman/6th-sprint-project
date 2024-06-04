package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"
	"time"

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
// ...

func getTasks(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(tasks) // сбор данных в json формат
	if err != nil {
		t := time.Now()
		errorInfo := fmt.Sprintln("Ошибка при сериализации данных в JSON формат", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusInternalServerError)                                        //выводы текста ошибок для упрощения определения на каком этапе получаем bad request
		fmt.Println("Ошибка при сериализации данных в JSON формат", t.Format("2006/01/02 15:04:05.00")) //все Println/Printf ошибок ниже придуманы для вывода информации на терминал
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(resp) // запись в тело ответа полученного ранее JSON
}

func postTask(w http.ResponseWriter, r *http.Request) {

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body) //читаем запрос, добавляем в буфер
	if err != nil {
		t := time.Now()
		errorInfo := fmt.Sprintln("Ошибка при чтении запроса", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusBadRequest)
		fmt.Println("Ошибка при чтении запроса", t.Format("2006/01/02 15:04:05.00"))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil { //десериализация данных из буфера
		t := time.Now()
		errorInfo := fmt.Sprintln("Ошибка при попытке десериализации", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusBadRequest)
		fmt.Println("Ошибка при попытке десериализации", t.Format("2006/01/02 15:04:05.00"))
		return
	}

	tasks[task.ID] = task //вносим в мапу данные по ключу равному значению поля ID структуры task (добавляем ключ и его значение)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func getTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id") // возврат параметра из запроса (Request)

	task, ok := tasks[id] //инициализация (+ проверка наличия в мапе) структуры по ключу(параметр полученный с помощью URLParam)
	if !ok {
		t := time.Now()
		errorInfo := fmt.Sprintln("Задачи с id =", id, "нет", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusBadRequest)
		fmt.Println("Задачи с id =", id, "нет", t.Format("2006/01/02 15:04:05.00"))
		return
	}

	resp, err := json.Marshal(task) //сериализация запрашиваемой структуры
	if err != nil {
		t := time.Now()
		errorInfo := fmt.Sprintln("Ошибка при сериализации данных в JSON формат", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusBadRequest)
		fmt.Println("Ошибка при сериализации данных в JSON формат", t.Format("2006/01/02 15:04:05.00"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id") // возврат параметра из запроса (Request)

	_, ok := tasks[id] //инициализация (+ проверка наличия в мапе) структуры по ключу(параметр полученный с помощью URLParam)
	if !ok {
		t := time.Now()
		errorInfo := fmt.Sprintln("Задачи с id =", id, "нет", t.Format("2006/01/02 15:04:05.00"))
		http.Error(w, errorInfo, http.StatusBadRequest)
		fmt.Println("Задачи с id =", id, "нет", t.Format("2006/01/02 15:04:05.00"))
		return
	}

	delete(tasks, id) // удаление пары ключ - значение из мапы по id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	r.Get("/tasks", getTasks) // регистрируем эндпоинт "/tasks" с методом Get

	r.Post("/tasks", postTask) // регистрируем эндпоинт "/tasks" с методом Post

	r.Get("/tasks/{id}", getTask) // регистрируем эндпоинт "/tasks{id}" с методом Get

	r.Delete("/tasks/{id}", deleteTask) //регистрируем эндпоинт "/tasks{id}" с методом Delete

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
