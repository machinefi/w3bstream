FROM golang:1.19

WORKDIR /go/bin

COPY ./build/srv-applet-mgr/srv-applet-mgr .
COPY ./build/srv-applet-mgr/openapi.json .
EXPOSE 8888

RUN echo $PATH
ENTRYPOINT ["/go/bin/srv-applet-mgr"]