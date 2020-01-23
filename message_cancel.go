package sonic

type ClientCancel struct{}

func (s ClientCancel) toProto() protoMessage {
	return protoMessage{
		MessageType: messageTypeCancel,
	}
}

func clientCancelFromProto(msg protoMessage) (Message, error) {
	return ClientCancel{}, nil
}

func (s ClientCancel) is_request() {
}
