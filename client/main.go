package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// curl "http://localhost:9000/cat-fact?name=Jamol"

var (
	drName     = "postgres"
	dbUrl      = "localhost"
	dbPort     = 5432
	dbName     = "sql"
	dbUser     = "postgres"
	dbPassword = "abdulaziz1221"
)

func DBConnect() *sqlx.DB {
	dbStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbUrl, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Open(drName, dbStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type Manager struct {
	DB *sqlx.DB
}

type Response struct {
	Message string `json:"message"`
}

var manager *Manager

func main() {
	manager = &Manager{DB: DBConnect()}

	port := ":9000"
	http.HandleFunc("/cat-fact", catFact)

	log.Println("The server is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func catFact(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("name")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	stringURL := "http://localhost:9001/cat-fact?name=" + str
	resp, err := client.Get(stringURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, resp.Status, http.StatusInternalServerError)
		return
	}

	var greetResp Response

	err = json.NewDecoder(resp.Body).Decode(&greetResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = manager.CreateCatFact(greetResp.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(greetResp)
}

func (m *Manager) CreateCatFact(message string) error {
	query := "INSERT INTO cat_fact (message) VALUES ($1)"
	_, err := m.DB.Exec(query, message)
	if err != nil {
		return err
	}
	return nil
}