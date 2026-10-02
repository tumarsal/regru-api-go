binary := "regru"
cmd := "./cmd/regru"
gobin := `go env GOBIN`
gopath_bin := `go env GOPATH` / "bin"

# путь установки: GOBIN или GOPATH/bin
install_dir := if gobin != "" { gobin } else { gopath_bin }

# собрать bin/regru
default: build

# показать доступные рецепты
help:
    @just --list

# собрать bin/regru
build:
    go build -o bin/{{ binary }} {{ cmd }}

# установить regru в {{ install_dir }}
install:
    go install {{ cmd }}

# go run ./cmd/regru …
run *ARGS:
    go run {{ cmd }} {{ ARGS }}

# unit-тесты
test:
    go test ./...

# проверить zone sync (экспорт YAML + dry-run); DOMAIN из аргумента или DOMAIN_NAME/.env
sync-test DOMAIN="" FILE="tmp-zone.yaml": build
    #!/usr/bin/env bash
    set -euo pipefail
    domain="{{ DOMAIN }}"
    if [[ -z "$domain" ]]; then
      domain="${DOMAIN_NAME:-}"
    fi
    if [[ -z "$domain" && -f .env ]]; then
      domain="$(grep -E '^DOMAIN_NAME=' .env | head -1 | cut -d= -f2- | tr -d '\r')"
    fi
    if [[ -z "$domain" ]]; then
      echo "укажите DOMAIN: just sync-test example.ru" >&2
      exit 1
    fi
    rm -f "{{ FILE }}"
    ./bin/{{ binary }} zone sync "$domain" -f "{{ FILE }}"
    test -f "{{ FILE }}"
    echo "--- содержимое {{ FILE }} ---"
    cat "{{ FILE }}"
    echo "--- dry-run (изменений быть не должно) ---"
    ./bin/{{ binary }} zone sync "$domain" -f "{{ FILE }}" --dry-run
    rm -f "{{ FILE }}"

# удалить bin/
clean:
    rm -rf bin/
