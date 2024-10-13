package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// LoginPage serves the login page.
func (hc *HandlerContext) LoginPage(rw http.ResponseWriter, r *http.Request) {
	returnTo := r.FormValue("returnTo")
	csrfToken := csrf.Token(r)
	user, _ := hc.UserFromSession(r)
	views.Page("Login", false, user, hc.pages.LoginPage(csrfToken, "", returnTo)).Render(rw)
}

// Authenticate authenticates a user.
func (hc *HandlerContext) Authenticate(rw http.ResponseWriter, r *http.Request) {
	returnTo := r.FormValue("returnTo")
	err := hc.sessions.RenewToken(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	hc.authprovider.Authenticate(rw, r, func(rw http.ResponseWriter, r *http.Request, authUser models.User) {
		userdata, err := json.Marshal(authUser)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}
		hc.sessions.Put(r.Context(), "user", string(userdata))
		slog.Info("Authenticated user", "user", authUser)
		slog.Info("Redirecting to", "returnTo", returnTo)
		http.Redirect(rw, r, returnTo, http.StatusSeeOther)
	})
}

// RequireAuth is a middleware that requires authentication.
func (hc *HandlerContext) RequireAuth(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		userdata := hc.sessions.GetString(r.Context(), "user")
		slog.Info("RequireAuth", "userdata", userdata)
		if userdata == "" {
			hxAwareRedirect(rw, r)
			return
		}
		var user models.User
		err := json.Unmarshal([]byte(userdata), &user)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}

func hxAwareRedirect(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		rw.Header().Set("HX-Redirect", "/login?returnTo="+r.Header.Get("HX-Current-URL"))
		rw.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(rw, r, "/login?returnTo="+r.URL.Path, http.StatusSeeOther)
}

// UserFromSession extracts the user from the session.
func (hc *HandlerContext) UserFromSession(r *http.Request) (models.User, error) {
	userdata := hc.sessions.GetString(r.Context(), "user")
	if userdata == "" {
		return models.User{}, nil
	}
	var user models.User
	err := json.Unmarshal([]byte(userdata), &user)
	if err != nil {
		slog.Error("Failed to unmarshal user data", "error", err)
		return models.User{}, err
	}
	return user, nil
}
