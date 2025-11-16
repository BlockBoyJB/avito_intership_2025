run:
	go run cmd/app/main.go

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down
