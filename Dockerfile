FROM golang:latest

WORKDIR /

COPY go.* .

RUN go mod download

COPY . .

RUN go build -o ./main ./cmd/api/

EXPOSE 22

CMD ["./main"]
