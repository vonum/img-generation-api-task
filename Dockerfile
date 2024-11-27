FROM golang:1.21.8-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o tuti-api cmd/api/server.go

EXPOSE 3000
CMD ["./tuti-api"]
