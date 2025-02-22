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

type Device struct {
	ID                         int            `json:"id"`
	Тип                        string         `json:"тип"`
	Название                   string         `json:"название"`
	Модель                     string         `json:"Модель"`
	Топливо                    string         `json:"Топливо"`
	Давление_атм               sql.NullString `json:"давление_атм"`
	Паропроизводительность     sql.NullString `json:"паропроизводительность, кг/ч"`
	Температура_пара           sql.NullString `json:"температура_пара"`
	КПД					       sql.NullString `json:"КПД"`
	Мощность           	       sql.NullString `json:"мощность"`
	Производство_пара_кг_ч     sql.NullString `json:"производство_пара, кг/ч"`
	Расход_газа                sql.NullString `json:"расход_газа"`
	Расход_дизеля		       sql.NullString `json:"расход_дизеля"`
	Расход_мазута              sql.NullString `json:"расход_мазута"`
	Расход_твердого_топлива    sql.NullString `json:"расход_твердого_топлива"`
	Вес_кг      			   sql.NullString `json:"вес, кг"`
	Расход_твердого_топлива3   sql.NullString `json:"расход_твердого_топлива3"`
	Расход_твердого_топлива4   sql.NullString `json:"расход_твердого_топлива4"`
}

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

	router.HandleFunc("/devices", getDevices).Methods("GET")
	router.HandleFunc("/devices/{id}", getDevice).Methods("GET")
	router.HandleFunc("/devices", createDevices).Methods("POST")
	router.HandleFunc("/devices/{id}", updateDevices).Methods("PUT")
	router.HandleFunc("/devices/{id}", deleteDevices).Methods("DELETE")

	
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")

	
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
			&d.КПД, 
			&d.Мощность,
			&d.Производство_пара_кг_ч,
			&d.Расход_газа,
			&d.Расход_дизеля,
			&d.Расход_мазута,
			&d.Расход_твердого_топлива, 
			&d.Вес_кг,
			&d.Расход_твердого_топлива3,
			&d.Расход_твердого_топлива4,
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
			&d.КПД, 
			&d.Мощность,
			&d.Производство_пара_кг_ч,
			&d.Расход_газа,
			&d.Расход_дизеля,
			&d.Расход_мазута,
			&d.Расход_твердого_топлива, 
			&d.Вес_кг,
			&d.Расход_твердого_топлива3,
			&d.Расход_твердого_топлива4,

	)
	if err != nil {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func createDevices(w http.ResponseWriter, r *http.Request) {
	var e Device
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := dbDevices.Exec(`INSERT INTO devices (
		"тип", "название", "Модель", "Топливо", "Давление, атм",
		"Паропроизводительность, кг/ч", "температура пара", "КПД", "мощность, кВт",
		"Производство пара, кг/ч", "Расход газа", "Расход дизеля", "Расход мазута",
		"Расход твердого топлива", "вес, кг", "Расход твердого топлива3", "Расход твердого топлива4"
	  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.Тип, e.Название, e.Модель, e.Топливо, e.Давление_атм, e.Паропроизводительность, e.Температура_пара, e.КПД,
		e.Мощность, e.Производство_пара_кг_ч, e.Расход_газа, e.Расход_дизеля, e.Расход_мазута, e.Расход_твердого_топлива, 
		e.Вес_кг, e.Расход_твердого_топлива3, e.Расход_твердого_топлива4)
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
func updateDevices(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"] // id из URL

    var e Device
    if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
        http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Выполняем UPDATE employees
    _, err := dbEmployees.Exec(`
        UPDATE devices
           SET "тип" = ?, "название" = ?, "Модель" = ?, "Топливо" = ?, "Давление, атм" = ?,
		   "Паропроизводительность, кг/ч" = ?, "температура пара" = ?, "КПД" = ?, "мощность, кВт" = ?,
		   "Производство пара, кг/ч" = ?, "Расход газа" = ?, "Расход дизеля" = ?, "Расход мазута" = ?,
		   "Расход твердого топлива" = ?, "вес, кг" = ?, "Расход твердого топлива3" = ?, "Расход твердого топлива4" = ?
         WHERE id = ?
    `, e.Тип, e.Название, e.Модель, e.Топливо, e.Давление_атм, e.Паропроизводительность, e.Температура_пара, e.КПД,
	e.Мощность, e.Производство_пара_кг_ч, e.Расход_газа, e.Расход_дизеля, e.Расход_мазута, e.Расход_твердого_топлива, 
	e.Вес_кг, e.Расход_твердого_топлива3, e.Расход_твердого_топлива4, id)
    if err != nil {
        http.Error(w, "Ошибка при обновлении записи: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(e)
}



// DELETE /employees/{id} – удаление сотрудника
func deleteDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := dbEmployees.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
