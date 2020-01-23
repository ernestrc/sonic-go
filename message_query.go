package sonic

import (
	"encoding/json"
	"errors"
)

type Query struct {
	AuthToken string
	Query     string
	Config    map[string]interface{}
}

func (s Query) toProto() protoMessage {
	payloadMap := make(map[string]interface{})
	payloadMap["config"] = s.Config
	if s.AuthToken != "" {
		payloadMap["auth"] = s.AuthToken
	}

	b, err := json.Marshal(payloadMap)
	if err != nil {
		panic(err)
	}

	return protoMessage{
		MessageType: messageTypeQuery,
		Variation:   s.Query,
		Payload:     b,
	}
}

func clientQueryFromProto(msg protoMessage) (Message, error) {
	query := Query{Query: msg.Variation}

	payloadMap := make(map[string]interface{})

	err := json.Unmarshal(msg.Payload, &payloadMap)
	if err != nil {
		return nil, err
	}

	auth, err := getStringInPayload(payloadMap, "auth")
	if err != nil {
		return nil, err
	}
	if auth != "" {
		query.AuthToken = auth
	}

	configIfc, ok := payloadMap["config"]
	if !ok {
		return nil, errors.New("missing config field in query")
	}

	config, ok := configIfc.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid type for query config")
	}
	query.Config = config

	return query, nil
}

func (s Query) is_request() {
}
