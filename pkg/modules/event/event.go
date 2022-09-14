// event management

package event

type Event interface {
	Meta() MetaData
	Raw() []byte
}

type MetaData struct {
	PublisherID string
	ProjectID   string
	AppletID    string // optional
}
