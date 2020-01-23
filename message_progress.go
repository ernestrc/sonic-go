package sonic

import (
	"encoding/json"
	"errors"
)

type StreamStatus int

const (
	StreamStatusQueued   StreamStatus = 0
	StreamStatusStarted               = 1
	StreamStatusRunning               = 2
	StreamStatusWaiting               = 3
	StreamStatusFinished              = 4
)

type StreamProgress struct {
	Status   StreamStatus
	Progress float64
	Total    float64
	Units    string
}

func (s StreamProgress) toProto() protoMessage {
	payloadMap := make(map[string]interface{})
	payloadMap["p"] = s.Progress
	payloadMap["s"] = s.Status
	payloadMap["t"] = s.Total
	payloadMap["u"] = s.Units

	b, err := json.Marshal(payloadMap)
	if err != nil {
		panic(err)
	}

	return protoMessage{
		MessageType: messageTypeProgress,
		Payload:     b,
	}
}

func streamProgressFromProto(msg protoMessage) (Message, error) {
	s := StreamProgress{}

	payloadMap := make(map[string]interface{})

	err := json.Unmarshal(msg.Payload, &payloadMap)
	if err != nil {
		return nil, err
	}

	s.Progress, err = getFloatInPayload(payloadMap, "p")
	if err != nil {
		return nil, err
	}
	if s.Progress == 0 {
		return nil, errors.New("missing progress field in progress message")
	}

	status, err := getFloatInPayload(payloadMap, "s")
	if err != nil {
		return nil, err
	}
	if status == 0 {
		return nil, errors.New("missing status field in progress message")
	}
	s.Status = StreamStatus(status)

	s.Total, err = getFloatInPayload(payloadMap, "t")
	if err != nil {
		return nil, err
	}

	s.Units, err = getStringInPayload(payloadMap, "u")
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s StreamProgress) is_response() {
}
