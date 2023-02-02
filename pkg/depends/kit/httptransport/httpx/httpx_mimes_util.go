package httpx

import "github.com/golang/protobuf/proto"

func NewApplicationProtobufWith(msg proto.Message) (*ApplicationProtobuf, error) {
	w := NewApplicationProtobuf()
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	return w, nil
}
