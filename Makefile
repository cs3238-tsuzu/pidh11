.PHONY: build cross copy
build:
	GOARCH=arm GOARM=6 GOOS=linux go build -o pidht11 .

cross:
	docker run -it --rm \
		-v $(shell pwd):/workspace \
		-w /workspace \
		-e CGO_ENABLED=1 \
		docker.elastic.co/beats-dev/golang-crossbuild:1.15.8-arm \
		--build-cmd "make build" \
		-p "linux/armv6"

copy: cross
	rsync -avh pidht11 pi@10.20.40.30:/home/pi/
