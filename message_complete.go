package sonic

import (
	"errors"
)

type StreamCompleted struct {
	Err error
}

func (s StreamCompleted) toProto() protoMessage {
	variation := s.Err.Error()

	return protoMessage{
		MessageType: messageTypeCompleted,
		Variation:   variation,
	}
}

func streamCompletedFromProto(msg protoMessage) (Message, error) {
	s := StreamCompleted{}
	if msg.Variation != "" {
		s.Err = errors.New(msg.Variation)
	}

	return s, nil
}

func (s StreamCompleted) is_response() {
}
