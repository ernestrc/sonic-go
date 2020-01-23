package sonic

import "encoding/json"

type EventType byte

const (
	EventTypeAuth      EventType = 'H'
	EventTypeCancel              = 'C'
	EventTypeStarted             = 'S'
	EventTypeQuery               = 'Q'
	EventTypeMeta                = 'T'
	EventTypeProgress            = 'P'
	EventTypeOutput              = 'O'
	EventTypeAck                 = 'A'
	EventTypeCompleted           = 'D'
)

type protoMessage struct {
	EventType EventType       `json:"e"`
	Variation string          `json:"v,omitempty"`
	Payload   json.RawMessage `json:"p,omitempty"`
}

type TypeHint struct {
	Name string
	Type json.Token
}

type TypeMetadata struct {
	TypesHint []TypeHint
}

func (a *TypeHint) UnmarshalJSON(b []byte) (err error) {
	kvs := make(map[string]json.Token)
	err = json.Unmarshal(b, &kvs)
	if err != nil {
		return
	}
	for k, v := range kvs {
		a.Name = k
		a.Type = v
		break
	}
	return
}

func (a TypeHint) MarshalJSON() ([]byte, error) {
	kvs := make(map[string]json.Token)
	kvs[a.Name] = a.Type
	return json.Marshal(kvs)
}

func (t *TypeMetadata) decode(msg protoMessage) (err error) {
	err = json.Unmarshal(msg.Payload, &t.TypesHint)
	return
}

func (t *TypeMetadata) encode(msg *protoMessage) {
	// err = json.Unmarshal(msg.Payload, &t.TypesHint)
}
