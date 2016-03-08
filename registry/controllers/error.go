package controllers

import (
	"encoding/json"
	"net/http"
)

func Error404(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"msg": "invalid url"}
	content, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(content))
}

func Error400(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"msg": "invalid request"}
	content, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(content))
}

func Error500(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"msg": "internal server error"}
	content, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(content))
}
