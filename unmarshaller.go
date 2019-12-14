package mono

import "encoding/json"

type unmarshaller struct {
}

func (u unmarshaller) Unmarshal(bts []byte, v interface{}) error {
	return json.Unmarshal(bts, v)
}
