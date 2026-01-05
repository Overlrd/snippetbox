all: run

run:
	go run ./cmd/web/ -addr=":8000"
