# env
## ARG DOCKER_REGISTRY=#{TODO}
## FROM ${DOCKER_REGISTRY}/#{TODO path to golang basic image} AS build-env
FROM golang:alpine AS build-env

RUN go version
RUN apk --no-cache add git
RUN apk --no-cache add openssh-client
RUN apk --no-cache add make
RUN apk --no-cache add ca-certificates

# add and start postgresql
RUN apk --no-cache add postgresql
ENV PG_DATA_DIR=/var/lib/postgresql/data
ENV PG_CFG_FILE=/var/lib/postgresql/data/postgresql.conf
CMD su postgres && mkdir ${PG_DATA_DIR} && chmod 0700 ${PG_DATA_DIR}
CMD su postgres && initdb -D ${PG_DATA_DIR}
RUN echo 'unix_socket_directories = /tmp' >> ${PG_CFG_FILE}
CMD su postgres && pg_ctl start -D ${PG_DATA_DIR}

# setup private depends
ARG REPO_ACESS_TOKEN=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
ARG REPO_HOST=github.com
ARG REPO_ORGNIZE=iotexproject
ARG REPO_NAME=w3bstream
ARG REPO_URL=${REPO_HOST}:${REPO_ORGNIZE}/${REPO_NAME}
# RUN git config --global url."git@${REPO_HOST}:".insteadOf "https://${REPO_ORGNIZE}:${REPO_ACESS_TOKEN}@${REPO_HOST}/"
# RUN git config --global url."ssh://git@${REPO_HOST}:".insteadOf "https://${REPO_ORGNIZE}:${REPO_ACESS_TOKEN}@${REPO_HOST}/"
ARG REPO=https://${REPO_ORGNIZE}:${REPO_ACESS_TOKEN}@${REPO_HOST}/${REPO_ORGNIZE}/${REPO_NAME}.git
RUN echo ${REPO}

WORKDIR /go/src
RUN git clone ${REPO}

WORKDIR /go/src/w3bstream
ENV GOPRIVATE=github.com/iotexproject/Bumblebee
RUN cd cmd/srv-applet-mgr && go build -mod vendor
RUN mkdir -pv build
RUN cp -r cmd/srv-applet-mgr/config build/
RUN mv cmd/srv-applet-mgr/srv-applet-mgr build/

EXPOSE 8888

ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}

WORKDIR /go/src/w3bstream/build

# ENTRYPOINT ["echo hello"]


