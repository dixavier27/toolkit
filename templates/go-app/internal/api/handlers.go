package api

import (
	"encoding/json"
	"net/http"
)

// NewMux monta o router usando o pattern routing do net/http (Go 1.22+).
func NewMux(version string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": version})
	})
	mux.HandleFunc("GET /hello/{nome}", helloHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	nome := r.PathValue("nome")
	if err := validateNome(nome); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"saudacao": "Olá, " + nome + "!"})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
