#!/usr/bin/env bash
# Sincroniza o snapshot do tupa-go em internal/tupavendor/source/.
#
# Uso:
#   scripts/sync-tupa-vendor.sh                 # clona HEAD da main
#   scripts/sync-tupa-vendor.sh v0.1.0          # clona uma tag/branch específica
#   TUPA_REPO=git@github.com:foo/tupa.git scripts/sync-tupa-vendor.sh
#
# Requer:
#   - git com acesso ao tupa-go (SSH ou HTTPS com token)
#   - bash, mktemp, cp, mv
#
# Saída esperada:
#   internal/tupavendor/source/{app,recurso,opcoes,ganchos,repositorio,contexto,erros}.gotxt
#   internal/tupavendor/source/LICENSE
#
# O script é idempotente: roda quantas vezes quiser para puxar a versão atual.

set -euo pipefail

REF="${1:-main}"
REPO="${TUPA_REPO:-git@github.com:dixavier27/tupa-go.git}"

# Resolve a raiz do projeto eco (diretório acima de scripts/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ECO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DEST="${ECO_ROOT}/internal/tupavendor/source"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "→ clonando ${REPO}@${REF} em ${TMP_DIR}"
git clone --quiet --depth 1 --branch "${REF}" "${REPO}" "${TMP_DIR}/tupa" \
  || (echo "fallback: clone full + checkout ${REF}" && \
      git clone --quiet "${REPO}" "${TMP_DIR}/tupa" && \
      git -C "${TMP_DIR}/tupa" checkout --quiet "${REF}")

# Quais arquivos copiar — toda raiz exceto cmd/, testes, docs, e configs do repo.
ARQUIVOS=(
  app.go
  recurso.go
  opcoes.go
  ganchos.go
  repositorio.go
  contexto.go
  erros.go
)

mkdir -p "${DEST}"

# Limpa snapshot anterior (apenas .gotxt e LICENSE, preserva README/notes).
echo "→ limpando snapshot anterior em ${DEST}"
find "${DEST}" -maxdepth 1 -type f \( -name '*.gotxt' -o -name 'LICENSE' \) -delete

# Copia LICENSE e cada arquivo .go renomeando para .gotxt.
echo "→ copiando snapshot novo"
if [ -f "${TMP_DIR}/tupa/LICENSE" ]; then
  cp "${TMP_DIR}/tupa/LICENSE" "${DEST}/LICENSE"
fi
for arq in "${ARQUIVOS[@]}"; do
  src="${TMP_DIR}/tupa/${arq}"
  if [ ! -f "${src}" ]; then
    echo "  ! ausente em tupa-go: ${arq}" >&2
    continue
  fi
  cp "${src}" "${DEST}/${arq}txt"
  echo "  + ${arq}txt"
done

# Captura o SHA snapshotado para rastreabilidade.
SHA="$(git -C "${TMP_DIR}/tupa" rev-parse HEAD)"
echo "${SHA}  ${REPO}  ${REF}" > "${DEST}/SOURCE_SHA"
echo "→ snapshot @ ${SHA}"

echo
echo "✓ vendor atualizado. Revise com:"
echo "    git -C ${ECO_ROOT} diff internal/tupavendor/source/"
echo
echo "Se OK, comite com algo como:"
echo "    git add internal/tupavendor/source"
echo "    git commit -m 'chore(vendor): sync tupa-go @ ${SHA:0:7}'"
