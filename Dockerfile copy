FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app/m3uproxy

RUN go build -o main .

CMD ["/app/m3uproxy/main"]
