FROM golang:alpine
MAINTAINER Sylvain Laurent

WORKDIR /
RUN apk add --no-cache curl && curl https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-64bit-static.tar.xz | tar xJv && \
    mv ffmpeg-4.0-64bit-static/ffmpeg /bin/ffmpeg

ENV GOBIN $GOPATH/bin
ENV PROJECT_DIR github.com/Magicking/secure-sunrise
ENV PROJECT_NAME secure-sunriset-server

ADD vendor /usr/local/go/src
ADD cmd /go/src/${PROJECT_DIR}/cmd
ADD models /go/src/${PROJECT_DIR}/models
ADD restapi /go/src/${PROJECT_DIR}/restapi
ADD internal /go/src/${PROJECT_DIR}/internal

ADD GeoLite2-City_20170606.mmdb /db.mmdb

WORKDIR /go/src/${PROJECT_DIR}

RUN go build -v -o /go/bin/main /go/src/${PROJECT_DIR}/cmd/${PROJECT_NAME}/main.go
ADD run.sh /go/src/${PROJECT_DIR}/
ADD get_sample.sh /
ENTRYPOINT /go/src/${PROJECT_DIR}/run.sh

EXPOSE 8090
