FROM alpine:3.10

RUN ls -l
RUN /m3uproxy.sh build

COPY bin/m3uproxy /usr/local/bin

CMD m3uproxy
