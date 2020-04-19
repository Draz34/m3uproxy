FROM webhippie/mariadb:latest

WORKDIR /srv/app
ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/bin/s6-svscan", "/etc/s6"]

ENV GOPATH /srv/app
ENV GO15VENDOREXPERIMENT 1

ENV PATH /srv/app/bin:/usr/local/go/bin:${PATH}

ENV GOLANG_VERSION 1.10.4
ENV GOLANG_TARBALL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz

RUN apk update && \
  apk upgrade && \
  apk add \
    build-base \
    git \
    git-lfs \
    mercurial \
    bzr \
    go && \
  export \
    GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
  curl -sLo - \
    ${GOLANG_TARBALL} | tar -xzf - -C /usr/local && \
  cd \
    /usr/local/go/src && \
  patch -p2 -i \
    /tmp/default-buildmode-pie.patch && \
  patch -p2 -i \
    /tmp/set-external-linker.patch && \
  bash \
    make.bash && \
  apk del \
    go && \
  rm -rf \
    /var/cache/apk/*



RUN mkdir /app

ADD . /app/

WORKDIR /app/m3uproxy

RUN go build -o main .


ARG VERSION
ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.version=$VERSION
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vcs-url="https://github.com/dockhippie/golang.git"
LABEL org.label-schema.name="Golang"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"