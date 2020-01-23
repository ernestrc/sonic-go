package sonic

import (
	"encoding/json"
	"errors"
)

type TypeHint struct {
	Name string
	Type interface{}
}

type StreamSchema []TypeHint

func (a *StreamSchema) UnmarshalJSON(b []byte) (err error) {
	kvs := make([][]interface{}, 0)
	err = json.Unmarshal(b, &kvs)
	if err != nil {
		return
	}
	*a = StreamSchema(make([]TypeHint, len(kvs)))

loop:
	for i, vmap := range kvs {
		for j, v := range vmap {
			switch j {
			case 0:
				k, ok := v.(string)
				if !ok {
					err = errors.New("first item in StreamSchema must be a string")
					return
				}
				[]TypeHint(*a)[i].Name = k
			case 1:
				[]TypeHint(*a)[i].Type = v
				break loop
			}
		}
	}
	return
}

func (a *StreamSchema) MarshalJSON() ([]byte, error) {
	kvs := make([][]interface{}, len([]TypeHint(*a)))
	for i, kv := range []TypeHint(*a) {
		kvs[i] = make([]interface{}, 2)
		kvs[i][0] = kv.Name
		kvs[i][1] = kv.Type
	}
	return json.Marshal(kvs)
}

func (s StreamSchema) toProto() protoMessage {
	msg := protoMessage{
		MessageType: messageTypeSchema,
	}

	var err error
	msg.Payload, err = json.Marshal(&s)
	if err != nil {
		panic(err)
	}
	return msg
}

func streamSchemaFromProto(in protoMessage) (Message, error) {
	s := StreamSchema{}
	err := json.Unmarshal(in.Payload, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s StreamSchema) is_response() {
}
