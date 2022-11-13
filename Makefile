generate:
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./pkg/ent ./internal/app/model --feature sql/modifier --feature sql/execquery --feature sql/versioned-migration

run:
	go run ./cmd/app/main.go

swag:
	swag init -g ./cmd/app/main.go
	
docker_postgres_dev:
	docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USERNAME=postgres -e TZ=Asia/Seoul -v ~/data/pgdata:/var/lib/postgresql/data -d postgres:latest

docker_postgres_test:
	docker-compose -f docker-compose.test.yml --env-file ./configs/.env.test up -d

apidoc:
	apidoc -c apidocs.json -i internal/app/delivery/http -o apidocs

mockery:
	mockery --output ./internal/app/mocks --recursive --all --dir ./internal/app && sh ./scripts/mock.sh

