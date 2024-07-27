swagger: 
	swag init --parseDependency --parseInternal --parseDepth 2 -g ../cmd/api/main.go -d handler -o ./cmd/api/docs

run:
	make swagger & go run cmd/api/main.go