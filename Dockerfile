FROM golang:latest

WORKDIR /

COPY go.* .

RUN go mod download

COPY . .

ENV TERM=xterm-256color

RUN go build -o ./main .

EXPOSE 22

CMD ["./main"]
