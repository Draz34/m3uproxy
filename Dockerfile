FROM alpine:3.10

RUN go install m3uproxy/main.go
RUN mv bin/main bin/m3uproxy

COPY bin/m3uproxy /usr/local/bin

CMD m3uproxy
