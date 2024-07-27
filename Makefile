swagger: 
	swag init -g ../cmd/api/main.go -d handler -o ./cmd/api/docs

run:
	make swagger & go run cmd/api/main.go