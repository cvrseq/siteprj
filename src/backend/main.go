package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Device описывает устройство с полями, соответствующими столбцам таблицы devices
type Device struct {
	ID                     int     `json:"id"`
	Тип                    string  `json:"тип"`
	Название               string  `json:"название"`
	Модель                 string  `json:"Модель"`
	Топливо                string  `json:"Топливо"`
	Давление_атм           float64 `json:"давление_атм"`
	Паропроизводительность float64 `json:"паропроизводительность, кг/ч"`
	Температура_пара       float64 `json:"температура_пара"`
	КПД                    float64 `json:"КПД"`
	Мощность_кВт           float64 `json:"мощность, кВт"`
	Производство_пара      float64 `json:"производство пара, кг/ч"`
	Расход_газа            float64 `json:"расход газа"`
	Расход_дизеля          float64 `json:"расход дизеля"`
	Расход_мазута          float64 `json:"расход мазута"`
	Расход_твердого_топлива float64 `json:"расход твердого топлива"`
	Вес_кг                 float64 `json:"вес, кг"`
	Расход_твердого_топлива3 float64 `json:"расход твердого топлива3"`
	Расход_твердого_топлива4 float64 `json:"расход твердого топлива4"`
}

var db *sql.DB

func main() {
	var err error
	// Открываем или создаем базу данных SQLite (файл mydatabase.db должен быть в той же директории)
	db, err = sql.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализируем маршрутизатор
	router := mux.NewRouter()

	// Определяем маршруты API
	router.HandleFunc("/devices", getDevices).Methods("GET")
	router.HandleFunc("/devices/{id}", getDevice).Methods("GET")
	router.HandleFunc("/devices", createDevice).Methods("POST")
	router.HandleFunc("/devices/{id}", updateDevice).Methods("PUT")
	router.HandleFunc("/devices/{id}", deleteDevice).Methods("DELETE")

	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// getDevices возвращает список всех устройств
func getDevices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM devices")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		err := rows.Scan(&d.ID, &d.Тип, &d.Название, &d.Модель, &d.Топливо, &d.Давление_атм,
			&d.Паропроизводительность, &d.Температура_пара, &d.КПД, &d.Мощность_кВт, &d.Производство_пара,
			&d.Расход_газа, &d.Расход_дизеля, &d.Расход_мазута, &d.Расход_твердого_топлива, &d.Вес_кг,
			&d.Расход_твердого_топлива3, &d.Расход_твердого_топлива4)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		devices = append(devices, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// getDevice возвращает устройство по его ID
func getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var d Device
	err := db.QueryRow("SELECT * FROM devices WHERE id = ?", id).
		Scan(&d.ID, &d.Тип, &d.Название, &d.Модель, &d.Топливо, &d.Давление_атм,
			&d.Паропроизводительность, &d.Температура_пара, &d.КПД, &d.Мощность_кВт, &d.Производство_пара,
			&d.Расход_газа, &d.Расход_дизеля, &d.Расход_мазута, &d.Расход_твердого_топлива, &d.Вес_кг,
			&d.Расход_твердого_топлива3, &d.Расход_твердого_топлива4)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// createDevice создает новое устройство
func createDevice(w http.ResponseWriter, r *http.Request) {
	var d Device
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`INSERT INTO devices (
        "тип", "название", "Модель", "Топливо", "Давление, атм",
        "Паропроизводительность, кг/ч", "температура пара", "КПД", "мощность, кВт",
        "Производство пара, кг/ч", "Расход газа", "Расход дизеля", "Расход мазута",
        "Расход твердого топлива", "вес, кг", "Расход твердого топлива3", "Расход твердого топлива4")
        VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		d.Тип, d.Название, d.Модель, d.Топливо, d.Давление_атм,
		d.Паропроизводительность, d.Температура_пара, d.КПД, d.Мощность_кВт,
		d.Производство_пара, d.Расход_газа, d.Расход_дизеля, d.Расход_мазута,
		d.Расход_твердого_топлива, d.Вес_кг, d.Расход_твердого_топлива3, d.Расход_твердого_топлива4)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastID, _ := result.LastInsertId()
	d.ID = int(lastID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// updateDevice обновляет устройство по ID
func updateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var d Device
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`UPDATE devices SET 
        "тип" = ?, "название" = ?, "Модель" = ?, "Топливо" = ?, "Давление, атм" = ?,
        "Паропроизводительность, кг/ч" = ?, "температура пара" = ?, "КПД" = ?, "мощность, кВт" = ?,
        "Производство пара, кг/ч" = ?, "Расход газа" = ?, "Расход дизеля" = ?, "Расход мазута" = ?,
        "Расход твердого топлива" = ?, "вес, кг" = ?, "Расход твердого топлива3" = ?, "Расход твердого топлива4" = ?
        WHERE id = ?`,
		d.Тип, d.Название, d.Модель, d.Топливо, d.Давление_атм,
		d.Паропроизводительность, d.Температура_пара, d.КПД, d.Мощность_кВт,
		d.Производство_пара, d.Расход_газа, d.Расход_дизеля, d.Расход_мазута,
		d.Расход_твердого_топлива, d.Вес_кг, d.Расход_твердого_топлива3, d.Расход_твердого_топлива4, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// deleteDevice удаляет устройство по ID
func deleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
