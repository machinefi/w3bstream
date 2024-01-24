# maybe use private hub?
#ARG DOCKER_REGISTRY=hub.docker.com
#ARG GO_VERSION=1.19
FROM golang:1.21 AS builder

# setup private pkg if needs
#ARG GITHUB_CI_TOKEN
#ARG GITHUB_HOST=github.com
#ARG GOPROXY=https://goproxy.cn,direct
#ENV GONOSUMDB=${GITHUB_HOST}/*
#ARG GOPRIVATE=${GITHUB_HOST}
#RUN git config --global url.https://github-ci-token:${GITHUB_CI_TOKEN}@${GITHUB_HOST}/.insteadOf https://${GITHUB_HOST}/

# FROM build-env AS builder

WORKDIR /go/src
COPY ./ ./

# build
#ARG COMMIT_SHA
RUN cd ./cmd/srv-applet-mgr && make target

# runtime
FROM golang:1.21 AS runtime

COPY --from=builder /go/src/build/srv-applet-mgr/srv-applet-mgr /go/bin/srv-applet-mgr
COPY --from=builder /go/src/build/srv-applet-mgr/openapi.json /go/bin/openapi.json
EXPOSE 8888

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/srv-applet-mgr"]
