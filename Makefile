-include .env
export

start:
	go run main.go

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	goose -dir migrations mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" up 

mysql:
	mysql -h 127.0.0.1 -u ${DB_USER} -p ${DB_NAME}
