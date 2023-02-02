package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang/protobuf/proto"

	testdatapb "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer/testdata"
)

var (
	url         = "http://127.0.0.1:8080/demo/binary/protobuf"
	contentType = "application/x-protobuf"
)

func main() {
	cli := &http.Client{}

	// body1, err := genProtobufBody(&testdatapb.Event{Header: &testdatapb.Header{EventId: "1"}})
	// checkError(err)
	// doRequestAndOutput(cli, body1)
	// expected:
	// {"event_id":"1","succeeded":true}

	body2, err := genProtobufBody(&testdatapb.Event{})
	checkError(err)
	doRequestAndOutput(cli, body2)
	// expected:
	// {"err_msg":"invalid event id"}
}

func doRequestAndOutput(cli *http.Client, body io.Reader) {
	rsp, err := cli.Post(url, contentType, body)
	checkError(err)

	data, err := io.ReadAll(rsp.Body)
	checkError(err)
	defer rsp.Body.Close()

	pbrsp := &testdatapb.HandleResult{}
	checkError(proto.Unmarshal(data, pbrsp))

	jsonrsp, err := json.Marshal(pbrsp)
	checkError(err)

	fmt.Println(string(jsonrsp))
}

func genProtobufBody(msg proto.Message) (*bytes.Buffer, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
