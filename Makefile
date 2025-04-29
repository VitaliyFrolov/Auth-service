build-dev:
	sudo docker-compose -f docker-compose.dev.yaml up -d --build

dev-down:
	sudo docker-compose -f docker-compose.dev.yaml down
