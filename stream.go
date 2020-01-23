package sonic

import "context"

type Stream interface {
	Receive(context.Context) (Message, error)
}
