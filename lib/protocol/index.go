package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"
)

func cleanJSONData(data []byte) []byte {
	// Remove null bytes and other control characters
	return bytes.Map(func(r rune) rune {
		if r == 0 || unicode.IsControl(r) {
			return -1
		}
		return r
	}, data)
}

type Request struct {
	Path    string `json:"path"`
	Numbers []int  `json:"numbers"`
}

type Response struct {
	Status   int    `json:"status"`
	Quantity int    `json:"quantity"`
	Number   int    `json:"number"`
	Body     []byte `json:"body"`
}

func (m *Request) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Error encoding request: %v", err)
	}
	return data, nil
}

func (m *Request) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("Error decoding request: %v", err)
	}
	return nil
}

func (m *Response) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Error encoding response: %v", err)
	}
	return data, nil
}

func (m *Response) Decode(data []byte) error {
	cleanData := cleanJSONData(data)
	err := json.Unmarshal(cleanData, m)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}
	return nil
}
