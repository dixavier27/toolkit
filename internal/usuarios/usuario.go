// Package usuarios é um CRUD de usuários que prova os módulos tupa (HTTP) e
// inmemdb (persistência) integrados. Persiste via repo.Repositorio[Usuario]
// e expõe endpoints HTTP via tupa. O hash de senha fica atrás de Hasher.
package usuarios

// Usuario é o recurso do domínio. SenhaCriptografada nunca é serializada em
// JSON (tag "-"); o que entra pela API é a senha em texto plano (EntradaCriar).
type Usuario struct {
	ID                 string `json:"id" bson:"_id,omitempty"`
	Nome               string `json:"nome" bson:"nome"`
	Sobrenome          string `json:"sobrenome" bson:"sobrenome"`
	Email              string `json:"email" bson:"email"`
	Whatsapp           string `json:"whatsapp" bson:"whatsapp"`
	SenhaCriptografada string `json:"-" bson:"senha"`
}

// EntradaCriar é o payload de criação/atualização: traz a senha em texto plano,
// que o serviço hasheia antes de persistir.
type EntradaCriar struct {
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
	Email     string `json:"email"`
	Whatsapp  string `json:"whatsapp"`
	Senha     string `json:"senha"`
}
