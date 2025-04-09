# Dockerfile for app1
FROM golang:1.20

WORKDIR /app
COPY appcontainer/main.go .
RUN go build -o app1 main.go

CMD ["./app1"]
