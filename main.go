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
	handlerContext := handlers.NewHandlerContext(storage.NewMemStore())
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
	_ = http.ListenAndServe("localhost:8080", mux)
}
