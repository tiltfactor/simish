package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tiltfactor/simish/domain"
	"github.com/tiltfactor/simish/impl"
	"github.com/tiltfactor/simish/test"
	"github.com/urfave/cli"
)

type inputData map[string]string

// App holds the db structure. Used for dep injection.
type App struct {
	db domain.InputOutputStore
}

type response struct {
	Input      string  `json:"input"`
	Response   string  `json:"response"`
	Match      string  `json:"match,omitempty"`
	Room       string  `json:"room"`
	Score      float64 `json:"score"`
	AiCol      int64   `json:"aiCol"`
	ResultType int64   `json:"resultType"`
}

func (r response) String() string {
	return fmt.Sprintf("Input: %s, Matched: %s, Score: %f, Response: %s, Room: %s",
		r.Input, r.Match, r.Score, r.Response, r.Room)
}

type rawData struct {
	Text  string `json:"text"`
	Reply string `json:"reply"`
}

// ResponseHandler handles the user request for a input output pair
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
		Input:      input[0],
		Response:   io.Output,
		Match:      io.Input,
		Room:       room[0],
		Score:      score,
		AiCol:      io.AiCol,
		ResultType: io.ResultType,
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

// DBConfig is used to import the db_cfg.json file
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

func startSimish(c *cli.Context) error {
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
	return nil
}

func readInput(prompt string) string {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt + " ")
	scan.Scan()
	return scan.Text()
}

func createCfg(c *cli.Context) error {

	fmt.Println("Follow along to create a db_cfg.")
	fmt.Println("Warning: This will overwrite existing db_cfg.json file")

	f, err := os.Create("./db_cfg.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	serverPort := readInput("Port to run Simish on:")
	user := readInput("MySQL database username:")
	pass := readInput("MySQL database password:")
	ip := readInput("MySQL database IP address:")
	port := readInput("MySQL database port:")
	db := readInput("MySQL database name:")

	cfg := DBConfig{
		ServerPort: serverPort,
		User:       user,
		Pass:       pass,
		IP:         ip,
		DB:         db,
		Port:       port,
	}

	cm, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println(string(cm))
	w := bufio.NewWriter(f)
	if _, err := w.Write(cm); err != nil {
		log.Fatal(err)
	}
	w.Flush()
	fmt.Println("Successfully created db_cfg.json")

	return nil
}

func runTest(c *cli.Context) {
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

	pairs := store.GetAllPairs(1)
	input := c.Args().Get(0)
	test.RunSoftMatch(input, pairs)
}

func main() {
	app := cli.NewApp()
	app.Name = "Simish"
	app.Usage = "Soft Matching Algorithm as a service"
	app.Version = "0.4.0"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Create db_cfg.json file",
			Action: createCfg,
		},
		{
			Name:   "start",
			Usage:  "Start the simish server",
			Action: startSimish,
		},
		{
			Name:   "test",
			Usage:  "Test the softmatch algorithm",
			Action: runTest,
		},
	}
	app.Run(os.Args)
}
