default:
	go build -o bin/guess-my-word

docker:
	docker-compose build
	docker-compose push
