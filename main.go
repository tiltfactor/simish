package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tiltfactor/simish/impl"
)

type inputData map[string]string

// App ...
type App struct {
	db *impl.SQLStore
}

type response struct {
	Input    string  `json:"input"`
	Response string  `json:"response"`
	Match    string  `json:"match,omitempty"`
	Room     string  `json:"room"`
	Score    float64 `json:"score"`
}

func (r response) String() string {
	return fmt.Sprintf("Input: %s, Matched: %s, Score: %f, Response: %s, Room: %s",
		r.Input, r.Match, r.Score, r.Response, r.Room)
}

type rawData struct {
	Text  string `json:"text"`
	Reply string `json:"reply"`
}

// ResponseHandler ...
func (a App) ResponseHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	input := vars["input"]
	room := vars["room"]

	if len(room) < 1 {
		http.Error(w, "Must provide room number", 400)
		return
	}

	roomNumber, err := strconv.ParseInt(room[0], 10, 64)
	if err != nil {
		http.Error(w, "Room number must be an int", 400)
		return
	}

	io, score := a.db.Response(input[0], roomNumber)

	resp := response{
		Input:    input[0],
		Response: io.Output,
		Match:    io.Input,
		Room:     room[0],
		Score:    score,
	}

	log.Println(resp)

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// DBConfig ...
type DBConfig struct {
	User       string `json:"username"`
	Pass       string `json:"password"`
	IP         string `json:"ip_addr"`
	Port       string `json:"db_port"`
	DB         string `json:"database"`
	ServerPort string `json:"server_port"`
}

func (cfg *DBConfig) connectionString() string {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.User, cfg.Pass, cfg.IP, cfg.DB)
	log.Println(connStr)
	return connStr
}

func main() {
	cfgFile, err := ioutil.ReadFile("./db_cfg.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &DBConfig{}
	if err := json.Unmarshal(cfgFile, cfg); err != nil {
		log.Fatal(err)
	}

	store, err := impl.NewSQLStore(cfg.connectionString())
	if err != nil {
		log.Fatal(err)
	}

	app := App{db: store}

	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/api/v1/response", app.ResponseHandler)

	// Bind to a port and pass our router in
	log.Println("Running on port:", cfg.ServerPort)

	logFile, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file")
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
