package sonic

type Message interface {
	toProto() protoMessage
}

type MessageRequest interface {
	Message

	is_request()
}

type MessageResponse interface {
	Message

	is_response()
}
