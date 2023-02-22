postgres:
	sudo docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:15-alpine
createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	sudo docker exec -it postgres15 dropdb simple_bank
migrateCreate:
	migrate create -ext=.sql -seq -dir=./migrations ${name}
migrateup:
	migrate -path=./migrations -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path=./migrations -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrateforce:
	migrate -path=./migrations -database="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" force ${version}

.PHONY: postgres createdb dropdb migrateCreate migrateup migratedown migrateforce