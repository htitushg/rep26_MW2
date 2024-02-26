package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type user struct {
	Name  string `json: "full_name"`
	Email string `json: "email_address"`
}

type apiHandler struct{}

// Ajout du 24/02/2024 19h20
var logs, _ = os.Create("logs/logs.log")
var jsonHandler = slog.NewJSONHandler(logs, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}).WithAttrs([]slog.Attr{slog.Int("Info", 13)})
var Logger = slog.New(jsonHandler)
var LogId = 0

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		end := time.Since(start)
		LogId++
		Logger.Info("Log() withLogger", slog.Int("reqId", LogId), slog.Duration("duration", end), slog.String("reqMethod", req.Method), slog.String("reqURL", req.URL.String()))
	})
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello from Home Handler %s", req.URL)
}

// http.Handler
func (apiHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u := user{Name: "Henry", Email: "htitushg@gmail.com"}
	rw.Header().Set("content-Type", "application/json")
	json.NewEncoder(rw).Encode(u)
}

func main() {
	mux := http.DefaultServeMux
	mux.Handle("/home", withLogger(http.HandlerFunc(homeHandler)))
	mux.Handle("/api", withLogger(apiHandler{}))
	http.ListenAndServe(":3000", mux)
}
