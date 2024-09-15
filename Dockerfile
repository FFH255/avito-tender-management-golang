FROM golang:1.22.5

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY ./src ./src
COPY ./migration ./migration

RUN go build -o main ./src/cmd/main.go

EXPOSE 8080

CMD ["./main"]