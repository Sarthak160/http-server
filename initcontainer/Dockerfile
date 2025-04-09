# Dockerfile for app2
FROM golang:1.20

WORKDIR /app
COPY initcontainer/main.go .
RUN go build -o app2 main.go

CMD ["./app2"]
