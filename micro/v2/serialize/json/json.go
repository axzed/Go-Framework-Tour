package json

import "encoding/json"

type Serializer struct {
}

func (s Serializer) Code() byte {
	return 1
}

func (s Serializer) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (s Serializer) Decode(data []byte, val any) error {
	return json.Unmarshal(data, val)
}
