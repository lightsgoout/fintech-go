FROM golang:latest
WORKDIR /fintech-go
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/fintech-go