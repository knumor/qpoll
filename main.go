package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/handlers"
	"github.com/knumor/qpoll/storage"
)

//go:embed public
var staticFiles embed.FS

var Env = "dev"

func main() {
	staticFs, _ := fs.Sub(staticFiles, "public")
	mux := http.NewServeMux()
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.FS(staticFs))))

	pollStorage := storage.NewSQLiteStore()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		slog.Info("Shutting down")
		pollStorage.Close()
		os.Exit(0)
	}()

	handlerContext := handlers.NewHandlerContext(pollStorage)
	setupRoutes(mux, handlerContext)

	secure := true
	if Env == "dev" {
		slog.Info("Running in dev mode")
		secure = false
	}
	CSRF := csrf.Protect(
		[]byte("32-long-byte-key-auth"),
		csrf.Path("/"),
		csrf.Secure(secure),
		csrf.FieldName("csrf_token"),
	)
	_ = http.ListenAndServe("0.0.0.0:8080", CSRF(mux))
}

func setupRoutes(mux *http.ServeMux, hc *handlers.HandlerContext) {
	mux.HandleFunc("GET /create/wordcloud", handlers.CreateWordCloudPage)
	mux.HandleFunc("POST /create/wordcloud", hc.CreateWordCloud)
	mux.HandleFunc("GET /wordcloud/{id}", hc.GetWordCloud)
	mux.HandleFunc("GET /create/multiple-choice", handlers.CreateMultipleChoicePage)
	mux.HandleFunc("POST /create/mc", hc.CreateMultipleChoice)
	mux.HandleFunc("GET /mc/{id}", hc.GetMultipleChoice)
	mux.HandleFunc("GET /create", handlers.CreatePage)
	mux.HandleFunc("GET /vote/{id}/", hc.VotePage)
	mux.HandleFunc("POST /vote", hc.VoteSubmit)
	mux.HandleFunc("POST /join", hc.JoinExistingPoll)
	mux.HandleFunc("GET /present/{id}", hc.PresentPoll)
	mux.HandleFunc("GET /qr/{id}", hc.GenQRForPoll)
	mux.HandleFunc("GET /", handlers.JoinPollPage)
}
