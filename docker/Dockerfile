FROM golang:1.12.0-alpine3.9

ENV SRC_DIR=/go/src/github.com/api

WORKDIR $SRC_DIR

ADD . $SRC_DIR

ENV GO111MODULE=on

RUN apk add --no-cache \
  alpine-sdk \
  git \
  openjdk8

RUN apk --no-cache add curl bash python
RUN curl -sSL https://sdk.cloud.google.com | bash

ENV PATH $PATH:/root/google-cloud-sdk/bin

RUN go get github.com/pilu/fresh
RUN go get github.com/jinzhu/gorm
RUN go get github.com/gin-gonic/gin


EXPOSE 8080

CMD ["fresh"]


# 0 gcloud config list
#    1 gcloud init
#    2 gcloud auth login
#    3 gcloud init
#    4 gcloud config list
#    5 gcloud projects list
#    6 gcloud config set project inbound-rock-249009
#    7 exit
#    8 gcloud init
#    9 gcloud projects list
#   10 gcloud components list
#   11 gcloud components update
#   12 gcloud components install beta cbt
#   13 gcloud config set project inbound-rock-249009
#   14 gcloud components list
#   15 gcloud beta emulators bigtable start
#   16 gcloud components install cloud-datastore-emulator
#   17 gcloud beta emulators datastore start
#   18 export DATASTORE_EMULATOR_HOST=localhost:8081
#   19 gcloud beta emulators datastore start
#   20 export DATASTORE_EMULATOR_HOST=localhost:8080
#   21 gcloud beta emulators datastore start
#   22 gcloud components install cloud-datastore-emulator
#   23 gcloud beta emulators datastore start
#   24 apt install openjdk-8-jre
#   25 apk add  openjdk-8-jre
#   26 apk add openjdk-8-jre
#   27 apk add openjdk-8-jre-headless
#   28 apk add openjdk8
#   29 gcloud beta emulators datastore start
#   30 history