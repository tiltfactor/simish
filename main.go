package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/tiltfactor/simish/domain"
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

type rawData struct {
	Text  string `json:"text"`
	Reply string `json:"reply"`
}

// ResponseHandler ...
func (a App) ResponseHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	input := vars["input"]
	room := vars["room"]
	io, _, score, _ := a.db.Response(input[0], room[0])

	resp := response{
		Input:    input[0],
		Response: io.Output,
		Match:    io.Input,
		Room:     room[0],
		Score:    score,
	}

	match := domain.NewMatch(input[0], io.Input, room[0])
	a.db.SaveMatch(match)

	data, _ := json.Marshal(resp)
	w.Write(data)
}

// UpvoteHandler ...
func (a App) UpvoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	input := vars["input"][0]
	match := vars["match"][0]
	room := vars["room"][0]
	a.db.Upvote(input, match, room)
}

// DownvoteHandler ...
func (a App) DownvoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	input := vars["input"][0]
	match := vars["match"][0]
	room := vars["room"][0]
	a.db.Downvote(input, match, room)
}

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

// DBConfig ...
type DBConfig struct {
	User string `json:"username"`
	Pass string `json:"password"`
	IP   string `json:"ip_addr"`
	Port string `json:"port"`
}

func (cfg *DBConfig) connectionString() string {
	return fmt.Sprintf("%s:%s@/%s", cfg.User, cfg.Pass, cfg.IP)
}

func main() {
	cfgFile, err := ioutil.ReadFile("./db_cfg.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &DBConfig{}
	if err := json.Unmarshal(cfgFile, cfg); err != nil {
		log.Fatal()
	}

	store, _ := impl.NewSQLStore(cfg.connectionString())
	app := App{
		db: store,
	}

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/api/v1/response", app.ResponseHandler)
	r.HandleFunc("/api/v1/upvote", app.UpvoteHandler)
	r.HandleFunc("/api/v1/downvote", app.DownvoteHandler)

	// Bind to a port and pass our router in
	http.ListenAndServe(":8765", r)
}
