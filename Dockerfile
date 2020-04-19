FROM webhippie/mariadb:latest

WORKDIR /srv/app
ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/bin/s6-svscan", "/etc/s6"]

ENV GOPATH /srv/app
ENV PATH /srv/app/bin:/usr/local/go/bin:${PATH}
ENV GO111MODULE auto

RUN apk add xz

RUN cd /tmp && \
  curl -sSL -o upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && \
  tar -xvf upx.tar.xz --strip-components=1

RUN apk update && \
  apk upgrade && \
  apk add gcc musl-dev openssl openssh-client make git git-lfs mercurial go protoc protobuf-dev binutils-gold && \
  export GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
  export GOOS="$(go env GOOS)" && \
  export GOARCH="$(go env GOARCH)" && \
  export GOHOSTOS="$(go env GOHOSTOS)" && \
  export GOHOSTARCH="$(go env GOHOSTARCH)" && \
  curl -sLo - https://golang.org/dl/go1.14.2.src.tar.gz | tar -xzf - -C /usr/local && \
  cd /usr/local/go/src && \
  bash make.bash && \
  rm -rf /usr/local/go/pkg/bootstrap /usr/local/go/pkg/obj && \
  apk del go && \
  rm -rf /var/cache/apk/*



RUN mkdir /app

ADD . /app/

WORKDIR /app/m3uproxy

RUN go build -o main .