default:
	buffalo build

arm:
	GOOS=linux GOARCH=arm GOARM=5 buffalo build

deploy: arm
	chmod +x ./bin/guess-my-word

	rsync ./bin/guess-my-word pi@10.0.1.2:/tmp/
	ssh pi@10.0.1.2 sudo mv /tmp/guess-my-word /usr/local/bin/guess-my-word
	ssh pi@10.0.1.2 sudo chown root:root /usr/local/bin/guess-my-word*

	nomad job run -address="http://10.0.1.2:4646" nomad.hcl
