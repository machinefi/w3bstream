package httpx

//go:generate go run internal/mimes/main.go
const (
	MIME_OCTET_STREAM      = "application/octet-stream"
	MIME_JSON              = "application/json"
	MIME_XML               = "application/xml"
	MIME_FORM_URLENCODED   = "application/x-www-form-urlencoded"
	MIME_MULTIPART_FORMDAT = "multipart/form-data"
	MIME_PROTOBUF          = "application/x-protobuf"
	MIME_MSGPACK           = "application/x-msgpack"
	MIME_PLAIN_TEXT        = "text/plain"
)
