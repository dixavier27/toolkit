package tupa

import (
	"encoding/json"
	"net/http"
)

// LerJSON decodifica o corpo da requisição em dst. Rejeita campos
// desconhecidos para falhar cedo em payloads malformados.
func LerJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// EscreverJSON serializa v como JSON com o status informado.
func EscreverJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(v)
}

// EscreverErro responde com um JSON {"erro": msg} e o status informado.
func EscreverErro(w http.ResponseWriter, status int, msg string) {
	_ = EscreverJSON(w, status, map[string]string{"erro": msg})
}
