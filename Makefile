run-receiver:
	reflex -s -r '\.go$$' go run -gcflags=-G=3 cmd/receiver/main.go

run-publisher:
	reflex -s -r '\.go$$' go run -gcflags=-G=3 cmd/publisher/main.go

play:
	go run -gcflags=-G=3 cmd/playground/main.go
