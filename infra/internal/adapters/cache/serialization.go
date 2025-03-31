package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func (c *ObjectCacher[T]) unmarshal(data []byte, out any) error {

	switch c.serializationType {

	case SerializationTypeJSON:
		return json.Unmarshal(data, out)

	case SerializationTypeGob:
		decoder := gob.NewDecoder(bytes.NewReader(data))
		return decoder.Decode(out)

	default:
		return nil
	}
}

func (c *ObjectCacher[T]) Marshal(in any) ([]byte, error) {

	switch c.serializationType {

	case SerializationTypeJSON:
		return json.Marshal(in)

	case SerializationTypeGob:
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(in); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil

	default:
		return nil, nil
	}
}
