.PHONY: run build test vet fmt seed tidy media docker-up docker-down

run:
	cd Server/MagicMovieStream && go run .

build:
	cd Server/MagicMovieStream && go build -o bin/magicstream .

test:
	cd Server/MagicMovieStream && go test ./...

vet:
	cd Server/MagicMovieStream && go vet ./...

fmt:
	cd Server/MagicMovieStream && gofmt -w .

tidy:
	cd Server/MagicMovieStream && go mod tidy

seed:
	cd Server/MagicMovieStream && go run ./cmd/seed

media:
	cd Server/MagicMovieStream && mkdir -p media

docker-up:
	docker compose up --build

docker-down:
	docker compose down
