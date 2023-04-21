#build backend
FROM golang:1.19 AS build-go

WORKDIR /w3bstream

COPY . /w3bstream/
RUN make



#run
FROM golang:1.19

WORKDIR /w3bstream

EXPOSE 8888

COPY --from=build-go /w3bstream/build/srv-applet-mgr /w3bstream/srv-applet-mgr

COPY cmd/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod a+x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
