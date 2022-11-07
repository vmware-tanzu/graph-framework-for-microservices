build:
	go build -o bin/calib main.go

docker_build:
	docker build -t nexus-calib:latest .
	docker tag nexus-calib:latest 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/nexus-calib:latest

run:
	go run main.go
