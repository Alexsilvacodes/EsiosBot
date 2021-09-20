FROM golang:1.14-alpine3.11

RUN apk add --no-cache git

WORKDIR /app/esios

COPY go.mod .

RUN go mod download

COPY . .

ENTRYPOINT ["go", "run", "main.go"]
