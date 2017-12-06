package main

import (
	"net/http"
	"github.com/markbates/goth/gothic"
	"encoding/json"
	"encoding/base64"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r  *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {

	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request)  {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	authCookieValue := make(map[string]interface{})
	authCookieValue["name"] = user.Name
	value, _ := json.Marshal(authCookieValue)
	valueBased := base64.StdEncoding.EncodeToString(value)

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: valueBased,
		Path:  "/"})

	redirectToChat(w)
}

func logoutHandler(w http.ResponseWriter, r *http.Request)  {
	gothic.Logout(w, r)
	redirectToHome(w)
}