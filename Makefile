lint:
	@golangci-lint run

prepare:
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -days 3650 \
      -nodes -keyout tls.key -out tls.crt -subj "/CN=localhost" \
      -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1"; \
      cp .env.example .env

test:
	@go test -race -count=1 ./...

build:
	@go build -gcflags="all=-N -l" -o http-analyzer

run:
	@go run .
