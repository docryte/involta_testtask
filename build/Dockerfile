FROM golang:latest

RUN mkdir /app

ADD .. /app

WORKDIR /app/cmd/main

EXPOSE 80/tcp

RUN go build -o main .

CMD ["/app/cmd/main/main"]