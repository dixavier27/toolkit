// Package gocheck inspeciona o ambiente local (toolchain Go e ferramentas
// auxiliares) e produz resultados estruturados para o comando doctor.
package gocheck

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// MinMajor / MinMinor é a versão mínima exigida do Go.
const (
	MinMajor = 1
	MinMinor = 26
)

// Status classifica o resultado de um check.
type Status int

const (
	StatusOK Status = iota
	StatusWarn
	StatusMissing
	StatusError
)

func (s Status) Symbol() string {
	switch s {
	case StatusOK:
		return "✓"
	case StatusWarn:
		return "!"
	case StatusMissing, StatusError:
		return "✗"
	}
	return "?"
}

// Check representa o resultado da inspeção de uma ferramenta ou recurso.
type Check struct {
	Name       string // "go", "golangci-lint", "air", "go.mod"
	Status     Status
	Version    string // ex: "1.26.3"
	Path       string // caminho do binário ou arquivo
	Message    string // descrição amigável do estado
	Suggestion string // próximo passo se não OK (comando de instalação ou link)
}

// OK indica se o check passou.
func (c Check) OK() bool { return c.Status == StatusOK }

// CheckGo verifica se o binário go está no PATH e satisfaz a versão mínima.
func CheckGo() Check {
	c := Check{Name: "go"}

	path, err := exec.LookPath("go")
	if err != nil {
		c.Status = StatusMissing
		c.Message = "go não encontrado no PATH"
		c.Suggestion = installInstructionsGo()
		return c
	}
	c.Path = path

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("falha ao executar `go version`: %v", err)
		return c
	}

	raw := strings.TrimSpace(string(out))
	re := regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)
	m := re.FindStringSubmatch(raw)
	if len(m) < 3 {
		c.Status = StatusError
		c.Message = fmt.Sprintf("não consegui interpretar a versão: %q", raw)
		return c
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	c.Version = strings.TrimPrefix(m[0], "go")

	if major < MinMajor || (major == MinMajor && minor < MinMinor) {
		c.Status = StatusWarn
		c.Message = fmt.Sprintf("Go %s instalado, mínimo exigido: %d.%d", c.Version, MinMajor, MinMinor)
		c.Suggestion = "Atualize em https://go.dev/dl/"
		return c
	}

	c.Status = StatusOK
	c.Message = fmt.Sprintf("Go %s", c.Version)
	return c
}

// CheckGolangciLint verifica se o golangci-lint está disponível.
func CheckGolangciLint() Check {
	c := Check{Name: "golangci-lint"}

	path, err := exec.LookPath("golangci-lint")
	if err != nil {
		c.Status = StatusMissing
		c.Message = "golangci-lint não encontrado no PATH"
		c.Suggestion = "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
		return c
	}
	c.Path = path

	out, err := exec.Command("golangci-lint", "version").CombinedOutput()
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("falha ao executar `golangci-lint version`: %v", err)
		return c
	}
	// Saída típica: "golangci-lint has version v1.59.1 built with go1.22.4 ..."
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	if m := re.FindStringSubmatch(string(out)); len(m) >= 2 {
		c.Version = m[1]
	}
	c.Status = StatusOK
	if c.Version != "" {
		c.Message = fmt.Sprintf("golangci-lint %s", c.Version)
	} else {
		c.Message = "golangci-lint instalado"
	}
	return c
}

// CheckAir verifica se o air (hot reload) está disponível.
func CheckAir() Check {
	c := Check{Name: "air"}

	path, err := exec.LookPath("air")
	if err != nil {
		c.Status = StatusMissing
		c.Message = "air não encontrado no PATH (opcional, usado para hot reload)"
		c.Suggestion = "go install github.com/air-verse/air@latest"
		return c
	}
	c.Path = path

	out, _ := exec.Command("air", "-v").CombinedOutput()
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	if m := re.FindStringSubmatch(string(out)); len(m) >= 2 {
		c.Version = m[1]
	}
	c.Status = StatusOK
	if c.Version != "" {
		c.Message = fmt.Sprintf("air %s", c.Version)
	} else {
		c.Message = "air instalado"
	}
	return c
}

// CheckCCompiler verifica se há um compilador C disponível no PATH, necessário
// para `go test -race` (o detector de race exige CGO_ENABLED=1, que por sua vez
// precisa de um compilador C). No Windows o Go não traz um compilador C, então
// é preciso instalar gcc separadamente. A ausência é apenas um alerta (Warn):
// o restante do fluxo funciona sem o detector de race.
func CheckCCompiler() Check {
	c := Check{Name: "compilador C (-race)"}

	for _, bin := range []string{"gcc", "clang", "cc"} {
		if path, err := exec.LookPath(bin); err == nil {
			c.Path = path
			c.Status = StatusOK
			c.Message = fmt.Sprintf("%s disponível — `go test -race` habilitado", bin)
			return c
		}
	}

	c.Status = StatusWarn
	c.Message = "nenhum compilador C no PATH — `go test -race` indisponível (exige CGO_ENABLED=1)"
	c.Suggestion = installInstructionsCC()
	return c
}

// CheckGoMod verifica se há um projeto Go no diretório atual (presença de go.mod).
// Retorna StatusOK quando há go.mod válido, StatusMissing quando não há projeto.
func CheckGoMod(dir string) Check {
	c := Check{Name: "go.mod"}

	modPath := filepath.Join(dir, "go.mod")
	info, err := os.Stat(modPath)
	if err != nil {
		c.Status = StatusMissing
		c.Message = "diretório atual não é um projeto Go (sem go.mod)"
		c.Suggestion = "Rode `eco new <nome>` para criar um, ou `go mod init <module>` em um diretório existente"
		return c
	}
	if info.IsDir() {
		c.Status = StatusError
		c.Message = "go.mod é um diretório, não um arquivo"
		return c
	}

	c.Path = modPath
	data, err := os.ReadFile(modPath)
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("falha ao ler go.mod: %v", err)
		return c
	}
	module := extractModule(string(data))
	if module == "" {
		c.Status = StatusWarn
		c.Message = "go.mod sem diretiva `module`"
		return c
	}
	c.Status = StatusOK
	c.Message = fmt.Sprintf("módulo %s", module)
	return c
}

// CheckAll executa todos os checks e devolve os resultados em ordem fixa.
// dir é usado para checks locais (go.mod); passe "" para pular esses.
func CheckAll(dir string) []Check {
	checks := []Check{
		CheckGo(),
		CheckCCompiler(),
		CheckGolangciLint(),
		CheckAir(),
	}
	if dir != "" {
		checks = append(checks, CheckGoMod(dir))
	}
	return checks
}

func extractModule(src string) string {
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}

func installInstructionsGo() string {
	switch runtime.GOOS {
	case "windows":
		return "winget install --id GoLang.Go  •  https://go.dev/dl/"
	case "darwin":
		return "brew install go  •  https://go.dev/dl/"
	case "linux":
		return "https://go.dev/dl/  (apt/dnf/pacman podem ter versão antiga)"
	default:
		return "https://go.dev/dl/"
	}
}

func installInstructionsCC() string {
	switch runtime.GOOS {
	case "windows":
		return "instale gcc (TDM-GCC https://jmeubank.github.io/tdm-gcc/ ou `winget install --id MSYS2.MSYS2`) e garanta CGO_ENABLED=1"
	case "darwin":
		return "xcode-select --install"
	case "linux":
		return "apt install build-essential  (ou dnf groupinstall 'Development Tools')"
	default:
		return "instale um compilador C (gcc/clang) e habilite CGO_ENABLED=1"
	}
}
