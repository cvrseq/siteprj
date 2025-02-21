package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Device описывает устройство
type Device struct {
	ID                     int            `json:"id"`
	Тип                    string         `json:"тип"`
	Название               string         `json:"название"`
	Модель                 string         `json:"Модель"`
	Топливо                string         `json:"Топливо"`
	Давление_атм           sql.NullString `json:"давление_атм"`
	Паропроизводительность sql.NullString `json:"паропроизводительность, кг/ч"`
	Температура_пара       sql.NullString `json:"температура_пара"`
	// Дополнительные поля можно добавить
}

// Employee описывает сотрудника для логина
type Employee struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Сырой пароль
	Role     string `json:"role"`
}

var (
	dbDevices   *sql.DB
	dbEmployees *sql.DB
)

func main() {
	var err error

	// Подключаемся к базе устройств
	dbDevices, err = sql.Open("sqlite3", "./db/mydatabase.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе устройств:", err)
	}
	defer dbDevices.Close()

	// Подключаемся к базе сотрудников
	dbEmployees, err = sql.Open("sqlite3", "./db/employees.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе сотрудников:", err)
	}
	defer dbEmployees.Close()

	router := mux.NewRouter()


	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	// CRUD-операции для устройств
	router.HandleFunc("/devices", getDevices).Methods("GET")
	router.HandleFunc("/devices/{id}", getDevice).Methods("GET")
	// Другие обработчики для устройств (create, update, delete) можно добавить

	// Эндпоинт логина: GET – отдает registration.html, POST – авторизация
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")

	// Раздача статических файлов из директории ../frontend/
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/")))

	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := dbEmployees.Query("SELECT id, username, password, role FROM employees")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		err := rows.Scan(&e.ID, &e.Username, &e.Password, &e.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, e)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}


func getEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var e Employee
    err := dbEmployees.QueryRow(
        "SELECT id, username, password, role FROM employees WHERE id = ?", id,
    ).Scan(&e.ID, &e.Username, &e.Password, &e.Role)

    if err != nil {
        http.Error(w, "Не удалось найти сотрудника с id "+id, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(e)
}


// POST /employees – добавление нового сотрудника
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := dbEmployees.Exec("INSERT INTO employees (username, password, role) VALUES (?, ?, ?)",
		e.Username, e.Password, e.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastID, _ := result.LastInsertId()
	e.ID = int(lastID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

// PUT /employees/{id} – обновление сотрудника
func updateEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"] // id из URL

    var e Employee
    if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
        http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Выполняем UPDATE employees
    _, err := dbEmployees.Exec(`
        UPDATE employees 
           SET username = ?, password = ?, role = ?
         WHERE id = ?
    `, e.Username, e.Password, e.Role, id)
    if err != nil {
        http.Error(w, "Ошибка при обновлении записи: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(e)
}



// DELETE /employees/{id} – удаление сотрудника
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := dbEmployees.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func getDevices(w http.ResponseWriter, r *http.Request) {
	rows, err := dbDevices.Query("SELECT * FROM devices")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		err := rows.Scan(
			&d.ID,
			&d.Тип,
			&d.Название,
			&d.Модель,
			&d.Топливо,
			&d.Давление_атм,
			&d.Паропроизводительность,
			&d.Температура_пара,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		devices = append(devices, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var d Device
	err := dbDevices.QueryRow("SELECT * FROM devices WHERE id = ?", id).Scan(
		&d.ID,
		&d.Тип,
		&d.Название,
		&d.Модель,
		&d.Топливо,
		&d.Давление_атм,
		&d.Паропроизводительность,
		&d.Температура_пара,
	)
	if err != nil {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Отдаем страницу регистрации (registration.html)
		http.ServeFile(w, r, "../frontend/pages/registration.html")
		return
	}

	// POST-запрос: ожидаем JSON с username и password
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Ищем сотрудника в базе сотрудников
	var storedPassword, role string
	err = dbEmployees.QueryRow(
		"SELECT password, role FROM employees WHERE username = ?",
		creds.Username,
	).Scan(&storedPassword, &role)
	if err != nil {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	// Сравниваем сырой пароль
	if storedPassword != creds.Password {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	// При успешной авторизации
	if role == "admin" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": "/pages/admin.html"})
	} else if role == "user" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": "/index.html"})
	} else {
		http.Error(w, "Неизвестная роль", http.StatusForbidden)
	}
}
