FROM golang:1.17-alpine3.14

RUN apk add --no-cache git

WORKDIR /app/esios

COPY go.mod .

RUN go mod download

COPY . .

ENTRYPOINT ["go", "run", "main.go"]
