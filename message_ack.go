package sonic

type ClientAck struct{}

func (s ClientAck) toProto() protoMessage {
	return protoMessage{
		MessageType: messageTypeAck,
	}
}

func clientAckFromProto(msg protoMessage) (Message, error) {
	return ClientAck{}, nil
}

func (s ClientAck) is_request() {
}
