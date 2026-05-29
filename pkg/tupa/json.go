package tupa

import (
	"encoding/json"
	"net/http"
)

// LimiteCorpo é o tamanho máximo aceito por LerJSON (1 MiB).
const LimiteCorpo = 1 << 20

// LerJSON decodifica o corpo da requisição em dst. Rejeita campos
// desconhecidos para falhar cedo em payloads malformados. Limita a leitura
// a LimiteCorpo bytes para prevenir DoS por payload gigante.
func LerJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, LimiteCorpo)
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
