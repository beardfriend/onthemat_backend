generate:
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./pkg/ent ./internal/app/model

run:
	go run ./cmd/app/main.go

docker_postgres_dev:
	docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USERNAME=postgres -e TZ=Asia/Seoul -v ~/data/pgdata:/var/lib/postgresql/data -d postgres:latest