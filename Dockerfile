FROM alpine:3.10

ADD . /app
RUN ls -l
RUN /app/m3uproxy.sh build

COPY bin/m3uproxy /usr/local/bin

CMD m3uproxy
