package sonic

import "encoding/json"

type ClientAuth struct {
	Key  string
	User string
}

func (s ClientAuth) toProto() protoMessage {
	payloadMap := make(map[string]string)
	payloadMap["user"] = s.User

	b, err := json.Marshal(payloadMap)
	if err != nil {
		panic(err)
	}

	return protoMessage{
		MessageType: messageTypeAuth,
		Variation:   s.Key,
		Payload:     b,
	}
}

func clientAuthFromProto(msg protoMessage) (Message, error) {
	auth := ClientAuth{Key: msg.Variation}
	payloadMap := make(map[string]string)

	err := json.Unmarshal(msg.Payload, &payloadMap)
	if err != nil {
		return nil, err
	}

	auth.User = payloadMap["user"]
	return auth, nil
}

func (s ClientAuth) is_request() {
}
