run-receiver:
	reflex -s -r '\.go$$' go run cmd/receiver/main.go

run-publisher:
	reflex -s -r '\.go$$' go run cmd/publisher/main.go
