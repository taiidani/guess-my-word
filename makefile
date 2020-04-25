default:
	go build -o bin/guess-my-word

arm:
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/guess-my-word

pack:
	rm -f pkged.go
	pkger

upload: pack arm
	chmod +x ./bin/guess-my-word

	rsync ./bin/guess-my-word pi@10.0.1.2:/tmp/
	ssh pi@10.0.1.2 sudo mv /tmp/guess-my-word /usr/local/bin/guess-my-word
	ssh pi@10.0.1.2 sudo chown root:root /usr/local/bin/guess-my-word*

# Run this when code has been changed
deploy: upload
	nomad job stop -address="http://10.0.1.2:4646" guess-my-word || true
	nomad job run -address="http://10.0.1.2:4646" nomad.hcl
