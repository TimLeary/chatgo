package main

import "net/http"

func redirectToHome(w http.ResponseWriter)  {
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func redirectToChat(w http.ResponseWriter)  {
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}