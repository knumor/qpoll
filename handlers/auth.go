package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/models"
)

// LoginPage serves the login page.
func (hc *HandlerContext) LoginPage(rw http.ResponseWriter, r *http.Request) {
	r.FormValue("returnTo")
	csrfToken := csrf.Token(r)
	_ = hc.pages.LoginPage(csrfToken, "", "/create").Render(rw)
}

// Authenticate authenticates a user.
func (hc *HandlerContext) Authenticate(rw http.ResponseWriter, r *http.Request) {
	returnTo := r.FormValue("returnTo")
	err := hc.sessions.RenewToken(r.Context())
	if err != nil {
	http.Error(rw, err.Error(), 500)
		return
	}
	dummyuser := models.User{
		Username: "dummyuser",
		ID:       "dummyid",
		Email:    "dummy@example.invalid",
		Name:     "Dummy User",
	}
	userdata, err := json.Marshal(dummyuser)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	hc.sessions.Put(r.Context(), "user", string(userdata))
	slog.Info("Authenticated user", "user", dummyuser)
	slog.Info("Redirecting to", "returnTo", returnTo)
	http.Redirect(rw, r, returnTo, http.StatusSeeOther)
	return
}
// 	csrfToken := csrf.Token(r)
// 	_ = views.LoginPage(csrfToken, "invalid credentials", returnTo).Render(rw)
// }
//

// RequireAuth is a middleware that requires authentication.
func (hc *HandlerContext) RequireAuth(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		userdata := hc.sessions.GetString(r.Context(), "user")
		slog.Info("RequireAuth", "userdata", userdata)
		if userdata == "" {
			http.Redirect(rw, r, "/login?returnTo="+r.URL.Path, http.StatusSeeOther)
			return
		}
		var user models.User
		err := json.Unmarshal([]byte(userdata), &user)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}
		hc.pages.AuthenticatedUser = user
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}
