package sonic

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type StreamOutput []interface{}

func (s StreamOutput) toProto() protoMessage {
	b, err := json.Marshal(&s)
	if err != nil {
		panic(err)
	}

	return protoMessage{
		MessageType: messageTypeOutput,
		Payload:     b,
	}
}

func streamOutputFromProto(msg protoMessage) (Message, error) {
	s := StreamOutput{}

	err := json.Unmarshal(msg.Payload, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s StreamOutput) is_response() {
}

func (s StreamOutput) Unmarshal(schema StreamSchema, out interface{}) {
	// FIXME perf sucks
	m := s.UnmarshalRaw(schema)
	err := mapstructure.Decode(m, out)
	if err != nil {
		panic(err)
	}
}

func (s StreamOutput) UnmarshalRaw(schema StreamSchema) map[string]interface{} {
	ret := make(map[string]interface{})
	for i, tp := range schema {
		ret[tp.Name] = s[i]
	}
	return ret
}
