default:
	go build -o bin/guess-my-word

docker:
	docker-compose build
	docker-compose push

deploy: docker
	kubectl apply -f dist/

pack:
	rm -f pkged.go
	pkger
