default:
	buffalo build

arm:
	GOOS=linux GOARCH=arm GOARM=5 buffalo build

pi: arm
	cd deploy && ansible-playbook deploy.yml -i "10.0.1.2,"
