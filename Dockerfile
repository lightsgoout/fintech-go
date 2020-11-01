FROM golang:latest
WORKDIR /fintech-go
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/fintech-go