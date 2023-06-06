FROM golang:1.20.4-alpine3.17

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 1321

CMD ["./main"]