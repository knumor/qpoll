package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/authproviders"
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

	secure := true
	baseURL := os.Getenv("BASE_URL")
	if Env == "dev" {
		slog.Info("Running in dev mode")
		secure = false
		baseURL = "http://localhost:8080"
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = secure

	feideAuthProvider := authproviders.NewFeideAuthProvider(baseURL + "/auth/callback")
	mux.Handle("GET /auth/callback", http.HandlerFunc(feideAuthProvider.AuthResponseHandler))

	handlerContext := handlers.NewHandlerContext(pollStorage, sessionManager, feideAuthProvider)
	setupRoutes(mux, handlerContext)

	CSRF := csrf.Protect(
		[]byte("2543b4efe309a66bcbf93a390086abf4"),
		csrf.Path("/"),
		csrf.Secure(secure),
		csrf.FieldName("csrf_token"),
	)
	_ = http.ListenAndServe(
		"0.0.0.0:8080",
		CSRF(
			sessionManager.LoadAndSave(
				handlerContext.EnsureClientID(mux),
			),
		),
	)
}

func setupRoutes(mux *http.ServeMux, hc *handlers.HandlerContext) {
	mux.Handle("GET /create/wordcloud", hc.RequireAuth(http.HandlerFunc(hc.CreateWordCloudPage)))
	mux.Handle("POST /create/wordcloud", hc.RequireAuth(http.HandlerFunc(hc.CreateWordCloud)))
	mux.Handle("GET /wordcloud/{id}", hc.RequireAuth(http.HandlerFunc(hc.GetWordCloud)))
	mux.Handle("GET /create/multiple-choice", hc.RequireAuth(http.HandlerFunc(hc.CreateMultipleChoicePage)))
	mux.Handle("POST /create/mc", hc.RequireAuth(http.HandlerFunc(hc.CreateMultipleChoice)))
	mux.Handle("GET /mc/{id}", hc.RequireAuth(http.HandlerFunc(hc.GetMultipleChoice)))
	mux.Handle("GET /create", hc.RequireAuth(http.HandlerFunc(hc.CreatePage)))
	mux.HandleFunc("GET /vote/{id}/", hc.VotePage)
	mux.HandleFunc("POST /vote", hc.VoteSubmit)
	mux.HandleFunc("POST /join", hc.JoinExistingPoll)
	mux.Handle("GET /present/{id}", hc.RequireAuth(http.HandlerFunc(hc.PresentPoll)))
	mux.Handle("GET /qr/{id}", hc.RequireAuth(http.HandlerFunc(hc.GenQRForPoll)))
	mux.Handle("GET /polls", hc.RequireAuth(http.HandlerFunc(hc.ListPollsPage)))
	mux.HandleFunc("GET /login", hc.LoginPage)
	mux.HandleFunc("POST /login", hc.Authenticate)
	mux.HandleFunc("GET /{$}", hc.JoinPollPage)
}
