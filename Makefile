build:
	go build .

run: swag build db
	go run .

docker.build:
	docker-compose build

docker.up: docker.build
	docker-compose up --build --remove-orphans

docker.down:
	docker-compose down 

db:
	docker start byte_me

swag:
	swag init

lint:
	golangci-lint run