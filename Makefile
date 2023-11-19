build:
	go build .

run: swag build db
	go run .

docker.build: swag
	docker-compose build

docker.up: docker.build
	docker-compose start database
	docker-compose start server

docker.down:
	docker-compose down 

db:
	docker start byte_me

swag:
	swag init