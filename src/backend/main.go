package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var store = sessions.NewCookieStore([]byte("секретный-ключ-в-продакшне"))



func stringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

type Device struct {
	ID            int            `json:"-"`
	Type          sql.NullString `json:"-"`
	Name          sql.NullString `json:"-"`
	Model         sql.NullString `json:"-"`
	Fuel          sql.NullString `json:"-"`
	PressureAtm   sql.NullString `json:"-"`
	SteamCapacity sql.NullString `json:"-"`
	SteamTemp     sql.NullString `json:"-"`
	Efficiency    sql.NullString `json:"-"`
	Power         sql.NullString `json:"-"`
	SteamProd     sql.NullString `json:"-"`
	GasCons       sql.NullString `json:"-"`
	DieselCons    sql.NullString `json:"-"`
	FuelOilCons   sql.NullString `json:"-"`
	SolidFuelCons sql.NullString `json:"-"`
	Weight        sql.NullString `json:"-"`
	Burner        sql.NullString `json:"-"`
	Mop           sql.NullString `json:"-"`
	Mtp           sql.NullString `json:"-"`
}

type deviceAlias struct {
	ID            int    `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Model         string `json:"model"`
	Fuel          string `json:"fuel"`
	PressureAtm   string `json:"pressure"`
	SteamCapacity string `json:"steam_capacity"`
	SteamTemp     string `json:"steam_temperature"`
	Efficiency    string `json:"efficiency"`
	Power         string `json:"power"`
	SteamProd     string `json:"steam_production"`
	GasCons       string `json:"gas_cons"`
	DieselCons    string `json:"diesel_cons"`
	FuelOilCons   string `json:"fuel_oil_cons"`
	SolidFuelCons string `json:"solid_fuel_cons"`
	Weight        string `json:"weight"`
	Burner        string `json:"burner"`
	Mop           string `json:"mop"`
	Mtp           string `json:"mpt"`
}

func (d Device) MarshalJSON() ([]byte, error) {
	alias := deviceAlias{
		ID:            d.ID,
		Type:          stringOrEmpty(d.Type),
		Name:          stringOrEmpty(d.Name),
		Model:         stringOrEmpty(d.Model),
		Fuel:          stringOrEmpty(d.Fuel),
		PressureAtm:   stringOrEmpty(d.PressureAtm),
		SteamCapacity: stringOrEmpty(d.SteamCapacity),
		SteamTemp:     stringOrEmpty(d.SteamTemp),
		Efficiency:    stringOrEmpty(d.Efficiency),
		Power:         stringOrEmpty(d.Power),
		SteamProd:     stringOrEmpty(d.SteamProd),
		GasCons:       stringOrEmpty(d.GasCons),
		DieselCons:    stringOrEmpty(d.DieselCons),
		FuelOilCons:   stringOrEmpty(d.FuelOilCons),
		SolidFuelCons: stringOrEmpty(d.SolidFuelCons),
		Weight:        stringOrEmpty(d.Weight),
		Burner:        stringOrEmpty(d.Burner),
		Mop:           stringOrEmpty(d.Mop),
		Mtp:           stringOrEmpty(d.Mtp),
	}
	return json.Marshal(alias)
}

type Employee struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var (
	dbDevices   *sql.DB
	dbEmployees *sql.DB
)

func main() {
	store = sessions.NewCookieStore([]byte("секретный-ключ-в-продакшне"))


store.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   86400 * 7, 
    HttpOnly: true,
    Secure:   false,      
    SameSite: http.SameSiteLaxMode,
}
	var err error

	dbDevices, err = sql.Open("sqlite3", "./db/mydatabase.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе устройств:", err)
	}
	defer dbDevices.Close()

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

	router.HandleFunc("/admin_employees.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "employees_panel", http.StatusMovedPermanently)
	}).Methods("GET")

	router.HandleFunc("/admin_employees", func(w http.ResponseWriter, r *http.Request) {
    	http.Redirect(w, r, "/employees_panel", http.StatusMovedPermanently)
	}).Methods("GET")

	router.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "user_panel", http.StatusMovedPermanently)
	}).Methods("GET")

	router.HandleFunc("/pages/admin.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "admin_panel", http.StatusMovedPermanently)
	}).Methods("GET")

	router.HandleFunc("/file_manager.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/files/", http.StatusMovedPermanently)
	}).Methods("GET")


	router.HandleFunc("/employees_panel", func(w http.ResponseWriter, r *http.Request) { 
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth { 
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}

		http.ServeFile(w, r, "../frontend/pages/admin_employees.html")
	}).Methods("GET")


	router.HandleFunc("/cache-clear", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth { 
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}).Methods("GET")

	router.HandleFunc("/upload", uploadHandler).Methods("POST")


	router.HandleFunc("/filemanager", func(w http.ResponseWriter, r *http.Request) { 
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth { 
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}

		http.ServeFile(w, r, "../frontend/pages/file_manager.html")
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		role, ok := session.Values["role"].(string)
		if !ok {
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}
		
		if role == "admin" {
			http.Redirect(w, r, "/admin_panel", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/user_panel", http.StatusSeeOther)
		}
	}).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)
		
		if ok && auth {
			role, ok := session.Values["role"].(string)
			if ok {
				if role == "admin" {
					http.Redirect(w, r, "/admin_panel", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/user_panel", http.StatusSeeOther)
				}
				return
			}
		}

		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "../frontend/pages/registration.html")
			return
		}
		
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		
		var storedPassword, role string
		err = dbEmployees.QueryRow("SELECT password, role FROM employees WHERE username = ?", creds.Username).
			Scan(&storedPassword, &role)
		if err != nil {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}
		
		if storedPassword != creds.Password {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}
		
		session.Values["authenticated"] = true
		session.Values["username"] = creds.Username
		session.Values["role"] = role
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Ошибка сохранения сессии: %v", err)
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}
		
		redirectURL := "/user_panel"
		if role == "admin" {
			redirectURL = "/admin_panel"
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": redirectURL})
	}).Methods("GET", "POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		
		session.Values["authenticated"] = false
		delete(session.Values, "username")
		delete(session.Values, "role")
		
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, "Ошибка при выходе из системы", http.StatusInternalServerError)
			return
		}
		
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}).Methods("GET")

	router.Use(corsMiddleware)

	go startFilebrowser()

	time.Sleep(2 * time.Second)

	
	

	router.HandleFunc("/user_panel", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)
		
		log.Printf("Доступ к /user_panel, auth=%v, ok=%v", auth, ok)
		
		if !ok || !auth {
			log.Printf("Перенаправление на /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		log.Printf("Отдаем файл index.html")
		http.ServeFile(w, r, "../frontend/index.html")
	}).Methods("GET")
	
	router.HandleFunc("/admin_panel", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)

		log.Printf("Доступ к /admin_panel, auth=%v, ok=%v", auth, ok)
		
		if !ok || !auth {
			log.Printf("Перенаправление на /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		role, ok := session.Values["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}
		
		http.ServeFile(w, r, "../frontend/pages/admin.html")
	}).Methods("GET")

	router.PathPrefix("/style/").Handler(http.FileServer(http.Dir("../frontend")))
	router.PathPrefix("/script/").Handler(http.FileServer(http.Dir("../frontend")))
	router.PathPrefix("/images/").Handler(http.FileServer(http.Dir("../frontend")))

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../frontend/"))))

	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
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
		if err := rows.Scan(&e.ID, &e.Username, &e.Password, &e.Role); err != nil {
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
	err := dbEmployees.QueryRow("SELECT id, username, password, role FROM employees WHERE id = ?", id).
		Scan(&e.ID, &e.Username, &e.Password, &e.Role)
	if err != nil {
		http.Error(w, "Не удалось найти сотрудника с id "+id, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

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

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err := dbEmployees.Exec("UPDATE employees SET username = ?, password = ?, role = ? WHERE id = ?",
		e.Username, e.Password, e.Role, id)
	if err != nil {
		http.Error(w, "Ошибка при обновлении записи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

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
	query := `
		SELECT
			id,
			COALESCE(type, '') AS type,
			COALESCE(name, '') AS name,
			COALESCE(model, '') AS model,
			COALESCE(fuel, '') AS fuel,
			COALESCE(pressure, '') AS pressure,
			COALESCE(steam_capacity, '') AS steam_capacity,
			COALESCE(steam_temperature, '') AS steam_temperature,
			COALESCE(efficiency, '') AS efficiency,
			COALESCE(power, '') AS power,
			COALESCE(steam_production, '') AS steam_production,
			COALESCE(gas_cons, '') AS gas_cons,
			COALESCE(diesel_cons, '') AS diesel_cons,
			COALESCE(fuel_oil_cons, '') AS fuel_oil_cons,
			COALESCE(solid_fuel_cons, '') AS solid_fuel_cons,
			COALESCE(weight, '') AS weight,
			COALESCE(burner, '') AS burner,
			COALESCE(modification_one_pump, '') AS mop,
			COALESCE(modification_two_pump, '') AS mpt
		FROM devices
	`
	rows, err := dbDevices.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(
			&d.ID,
			&d.Type,
			&d.Name,
			&d.Model,
			&d.Fuel,
			&d.PressureAtm,
			&d.SteamCapacity,
			&d.SteamTemp,
			&d.Efficiency,
			&d.Power,
			&d.SteamProd,
			&d.GasCons,
			&d.DieselCons,
			&d.FuelOilCons,
			&d.SolidFuelCons,
			&d.Weight,
			&d.Burner,
			&d.Mop,
			&d.Mtp,
		); err != nil {
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

	query := `
		SELECT
			id,
			COALESCE(type, '') AS type,
			COALESCE(name, '') AS name,
			COALESCE(model, '') AS model,
			COALESCE(fuel, '') AS fuel,
			COALESCE(pressure, '') AS pressure,
			COALESCE(steam_capacity, '') AS steam_capacity,
			COALESCE(steam_temperature, '') AS steam_temperature,
			COALESCE(efficiency, '') AS efficiency,
			COALESCE(power, '') AS power,
			COALESCE(steam_production, '') AS steam_production,
			COALESCE(gas_cons, '') AS gas_cons,
			COALESCE(diesel_cons, '') AS diesel_cons,
			COALESCE(fuel_oil_cons, '') AS fuel_oil_cons,
			COALESCE(solid_fuel_cons, '') AS solid_fuel_cons,
			COALESCE(weight, '') AS weight,
			COALESCE(burner, '') AS burner,
			COALESCE(modification_one_pump, '') AS mop,
			COALESCE(modification_two_pump, '') AS mpt
		FROM devices
		WHERE id = ?
	`
	var d Device
	err := dbDevices.QueryRow(query, id).Scan(
		&d.ID,
		&d.Type,
		&d.Name,
		&d.Model,
		&d.Fuel,
		&d.PressureAtm,
		&d.SteamCapacity,
		&d.SteamTemp,
		&d.Efficiency,
		&d.Power,
		&d.SteamProd,
		&d.GasCons,
		&d.DieselCons,
		&d.FuelOilCons,
		&d.SolidFuelCons,
		&d.Weight,
		&d.Burner,
		&d.Mop,
		&d.Mtp,
	)
	if err != nil {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func deleteDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := dbDevices.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	r.ParseMultipartForm(30 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	targetPath := filepath.Join(uploadDir, header.Filename)
	outFile, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл сохранен в %s", targetPath)
}

func startFilebrowser() { 
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) { 
		os.MkdirAll(uploadDir, 0755)
	}

	cmd := exec.Command("filebrowser", "-r", uploadDir, "-p", "8081", "-a", "0.0.0.0")
	if err := cmd.Start(); err != nil { 
		log.Fatalf("Ошибка запуска filebrowser: %v", err)
	}
	log.Printf("filebrowser запущен с PID %d", cmd.Process.Pid)
}




func updateDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var d Device
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	query := `
		UPDATE devices
		   SET type = ?, name = ?, model = ?, fuel = ?, pressure = ?,
		       steam_capacity = ?, steam_temperature = ?, efficiency = ?, power = ?,
		       steam_production = ?, gas_cons = ?, diesel_cons = ?, fuel_oil_cons = ?,
		       solid_fuel_cons = ?, weight = ?, burner = ?, modification_one_pump = ?,
		       modification_two_pump = ?
		 WHERE id = ?`
	_, err := dbDevices.Exec(query,
		d.Type, d.Name, d.Model, d.Fuel, d.PressureAtm, d.SteamCapacity, d.SteamTemp, d.Efficiency,
		d.Power, d.SteamProd, d.GasCons, d.DieselCons, d.FuelOilCons, d.SolidFuelCons,
		d.Weight, d.Burner, d.Mop, d.Mtp, id)
	if err != nil {
		http.Error(w, "Ошибка при обновлении записи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func createDevices(w http.ResponseWriter, r *http.Request) {
	var d Device
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `
		INSERT INTO devices (
			type, name, model, fuel, pressure, steam_capacity,
			steam_temperature, efficiency, power, steam_production,
			gas_cons, diesel_cons, fuel_oil_cons, solid_fuel_cons,
			weight, burner, modification_one_pump, modification_two_pump
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	result, err := dbDevices.Exec(query,
		d.Type, d.Name, d.Model, d.Fuel, d.PressureAtm, d.SteamCapacity, d.SteamTemp, d.Efficiency,
		d.Power, d.SteamProd, d.GasCons, d.DieselCons, d.FuelOilCons, d.SolidFuelCons,
		d.Weight, d.Burner, d.Mop, d.Mtp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastID, _ := result.LastInsertId()
	d.ID = int(lastID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

