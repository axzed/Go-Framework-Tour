package json

import "encoding/json"

type Serializer struct {

}

func (s Serializer) Code() byte {
	return 2
}

func (s Serializer) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (s Serializer) Decode(data []byte, val any) error {
	// if len(data) == 0 {
	// 	return nil
	// }
	return json.Unmarshal(data, val)
}

