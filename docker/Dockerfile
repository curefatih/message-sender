# Use the official Golang image as a base image
FROM golang:1.22.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal --parseGoList=false --parseDepth 2 -g ../cmd/api/main.go -d handler -o ./cmd/api/docs

ENV CGO_ENABLED=0

RUN go build -o /app/main cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/config.yml .

EXPOSE 8080

RUN ls

ENV GIN_MODE=release

CMD ["/root/main"]