package protocol

import "strings"

type Request struct {
	Path string
	Body []byte
}

type Response struct {
	Body []byte
}

func (m *Request) Encode() []byte {
	return []byte(m.Path + "\n" + string(m.Body))
}

func (m *Request) Decode(data []byte) {
	parts := strings.Split(string(data), "\n")
	m.Path = parts[0]
	m.Body = []byte(parts[1])
}

func (m *Response) Encode() []byte {
	return []byte(string(m.Body))
}

func (m *Response) Decode(data []byte) {
	m.Body = data
}
