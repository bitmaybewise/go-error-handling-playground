default:
	go run -gcflags=-G=3 cmd/receiver/main.go
	
run-receiver:
	go run -gcflags=-G=3 cmd/receiver/main.go

run-publisher:
	go run -gcflags=-G=3 cmd/publisher/main.go

play:
	go run -gcflags=-G=3 cmd/playground/main.go
