package oas

func NewServer(url string) *Server {
	s := &Server{}
	s.URL = url
	return s
}

type WithServers struct {
	Servers []*Server `json:"servers,omitempty"`
}

func (o *WithServers) AddServer(s *Server) {
	if s == nil {
		return
	}
	o.Servers = append(o.Servers, s)
}

type Server struct {
	ServerObject
	SpecExtensions
}

func (i Server) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.ServerObject, i.SpecExtensions)
}

func (i *Server) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.ServerObject, &i.SpecExtensions)
}

type ServerObject struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

func (o *ServerObject) AddVariable(key string, v *ServerVariable) {
	if v == nil {
		return
	}
	if o.Variables == nil {
		o.Variables = make(map[string]*ServerVariable)
	}
	o.Variables[key] = v
}

func NewServerVariable(defaultValue string) *ServerVariable {
	return &ServerVariable{
		ServerVariableObject: ServerVariableObject{
			Default: defaultValue,
		},
	}
}

type ServerVariable struct {
	ServerVariableObject
	SpecExtensions
}

func (i ServerVariable) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.ServerVariableObject, i.SpecExtensions)
}

func (i *ServerVariable) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.ServerVariableObject, &i.SpecExtensions)
}

type ServerVariableObject struct {
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}
