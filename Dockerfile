FROM golang:alpine AS build
WORKDIR /go/src/myapp
COPY . .
RUN go build -o /go/bin/myapp main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/myapp /go/bin/myapp
ENTRYPOINT ["/go/bin/myapp"]
