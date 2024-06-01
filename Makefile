test:
	docker-compose up -d
	sleep 2 # sleep 2s to ensure database is up and ready to accept connections
	go test ./...
	docker-compose down -v