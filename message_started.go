package sonic

type StreamStarted struct{}

func (s StreamStarted) toProto() protoMessage {
	return protoMessage{
		MessageType: messageTypeStarted,
	}
}

func streamStartedFromProto(msg protoMessage) (Message, error) {
	return StreamStarted{}, nil
}

func (s StreamStarted) is_response() {
}
