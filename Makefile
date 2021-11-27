test:
	go test --v
clean:
	rm registry
build:
	go build -o registry main.go
run:
	make build
	PORT=:9999 ./registry
	# curl http://localhost:9999/details
	# curl http://localhost:9999/shutdown