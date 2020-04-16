FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

RUN ls -l

CMD ["/app/main"]
