package authproviders

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/knumor/qpoll/handlers"
	"github.com/knumor/qpoll/models"
	"golang.org/x/oauth2"
)

// FeideAuthProvider is an authentication provider for Feide.
type FeideAuthProvider struct {
	oidcProvider *oidc.Provider
	config       oauth2.Config
	authStates   map[string]*authState
	mu           sync.Mutex
}

type authState struct {
	ctx          context.Context
	pkceVerifier string
	callback     handlers.AuthCallbackFunc
}

// NewFeideAuthProvider creates a new FeideAuthProvider.
func NewFeideAuthProvider(redirectURL string) *FeideAuthProvider {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://auth.dataporten.no")
	if err != nil {
		panic(err)
	}

	config := oauth2.Config{
		ClientID:     "83c40709-8aa6-44c7-93b0-f2ebb2255b26",
		ClientSecret: os.Getenv("FEIDE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "userinfo-name"},
	}

	return &FeideAuthProvider{
		oidcProvider: provider,
		config:       config,
		authStates:   make(map[string]*authState),
	}
}

// Authenticate authenticates a user via redirect to Feide.
func (fa *FeideAuthProvider) Authenticate(rw http.ResponseWriter, r *http.Request, callback handlers.AuthCallbackFunc) {
	fa.mu.Lock()
	defer fa.mu.Unlock()
	state := generateState(16)
	authState := &authState{
		ctx:          context.Background(),
		callback:     callback,
		pkceVerifier: oauth2.GenerateVerifier(),
	}
	fa.authStates[state] = authState
	http.Redirect(rw, r, fa.config.AuthCodeURL(state, oauth2.S256ChallengeOption(authState.pkceVerifier)), http.StatusFound)
}

// AuthResponseHandler handles the callback from Feide
func (fa *FeideAuthProvider) AuthResponseHandler(rw http.ResponseWriter, r *http.Request) {
	fa.mu.Lock()
	defer fa.mu.Unlock()
	stateParam := r.URL.Query().Get("state")
	authState, ok := fa.authStates[stateParam]
	if !ok {
		slog.Error("Invalid state")
		http.Error(rw, "Failed to perform authentication", http.StatusInternalServerError)
	}
	token, err := fa.config.Exchange(authState.ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(authState.pkceVerifier))
	if err != nil {
		slog.Error("Failed to exchange token", "error", err)
		http.Error(rw, "Failed to perform authentication", http.StatusInternalServerError)
	}
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		slog.Error("No ID token in token response")
		http.Error(rw, "Failed to perform authentication", http.StatusInternalServerError)
	}

	// Parse and verify ID Token payload.
	verifier := fa.oidcProvider.Verifier(&oidc.Config{ClientID: fa.config.ClientID})
	idToken, err := verifier.Verify(authState.ctx, rawIDToken)
	if err != nil {
		slog.Error("Failed to verify ID token", "error", err)
		http.Error(rw, "Failed to perform authentication", http.StatusInternalServerError)
	}

	// Extract custom claims
	var claims struct {
		Name  string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to extract claims", "error", err)
		http.Error(rw, "Failed to perform authentication", http.StatusInternalServerError)
	}
	delete(fa.authStates, stateParam)
	authState.callback(rw, r, models.User{
		Username: idToken.Subject,
		Name:     claims.Name,
	})
}

func generateState(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)

	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
