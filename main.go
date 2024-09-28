package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/knumor/qpoll/handlers"
	"github.com/knumor/qpoll/storage"
)

//go:embed public
var staticFiles embed.FS

func main() {
	staticFs, _ := fs.Sub(staticFiles, "public")
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /create/wordcloud", handlers.CreateWordCloudPage)
	mux.HandleFunc("POST /create/wordcloud", handlerContext.CreateWordCloud)
	mux.HandleFunc("GET /wordcloud/{id}", handlerContext.GetWordCloud)
	mux.HandleFunc("GET /create", handlers.CreatePage)
	mux.HandleFunc("GET /vote/{id}/", handlerContext.VotePage)
	mux.HandleFunc("POST /vote", handlerContext.VoteSubmit)
	mux.HandleFunc("POST /join", handlerContext.JoinExistingPoll)
	mux.HandleFunc("GET /present/{id}", handlerContext.PresentPoll)
	mux.HandleFunc("GET /qr/{id}", handlerContext.GenQRForPoll)
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.FS(staticFs))))
	mux.HandleFunc("GET /", handlers.JoinPollPage)
	_ = http.ListenAndServe("0.0.0.0:8080", mux)
}
