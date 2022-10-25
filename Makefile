generate:
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./pkg/ent ./internal/app/model

run:
	go run ./cmd/app/main.go