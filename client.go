package sonic

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	errCanceledCtx = errors.New("context canceled while waiting for response")
)

type ClientOption func(*Client) *Client

type Client struct {
	address       string
	mu            sync.Mutex
	clientMu      sync.Mutex
	conn          net.Conn
	c             chan error
	clientTimeout time.Duration
}

func NewClient(address string, clientOpts ...ClientOption) (client *Client) {
	client = new(Client)
	client.c = make(chan error)
	client.address = address

	defaultOpts := []ClientOption{
		WithRouteTimeout(5 * time.Second),
	}
	clientOpts = append(defaultOpts, clientOpts...)

	for _, o := range clientOpts {
		client = o(client)
	}

	return client
}

func WithRouteTimeout(d time.Duration) ClientOption {
	return func(c *Client) *Client {
		c.clientTimeout = d
		return c
	}
}

func (s *Client) closeConn(reason error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return
	}

	fields := log.Fields{
		"reason":    reason,
		keyCallType: "ClientCloseConnection",
		keyStep:     valueStepAttempt,
	}
	log.WithFields(fields).Trace()
	if closeErr := s.conn.Close(); closeErr != nil {
		fields[keyStep] = valueStepFailure
		fields[keyError] = closeErr
	} else {
		fields[keyStep] = valueStepSuccess
	}
	log.WithFields(fields).Warning()
	s.conn = nil
}

func (s *Client) setupConn(ctx context.Context) (
	conn net.Conn, err error,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		fields := log.Fields{
			keyCallType: "ClientSetupConnection",
			keyStep:     valueStepAttempt,
			"addr":      s.address,
		}
		log.WithFields(fields).Trace()
		dialer := net.Dialer{}
		s.conn, err = dialer.DialContext(ctx, "tcp", s.address)
		if err != nil {
			fields[keyStep] = valueStepFailure
			fields[keyError] = err
			log.WithFields(fields).Warning()
			return
		}
		fields[keyStep] = valueStepSuccess
		fields["local-addr"] = s.conn.LocalAddr().String()
		log.WithFields(fields).Debug()
	}

	conn = s.conn
	return
}

func writeMessage(writer io.Writer, msg Message) error {
	err := doWriteBuffer(writer, msg.toProto().encode())
	return err
}

func readMessage(reader io.Reader) (msg Message, err error) {
	var buf []byte
	buf, err = doReadBuffer(reader)
	if err != nil {
		return
	}
	var protoMsg protoMessage
	msg, err = protoMsg.decode(buf)
	return
}

func (s *Client) Stream(ctx context.Context, req MessageRequest) (
	<-chan MessageResponse, error,
) {
	s.clientMu.Lock()

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, s.clientTimeout)
	defer cancel()

	tx := make(chan MessageResponse)

	go func() {
		// FIXME multiplex connections then
		// lock can be released when this function returns
		defer s.clientMu.Unlock()

		var err error
		var conn net.Conn
		var msg Message

		conn, err = s.setupConn(ctx)
		if err != nil {
			s.closeConn(err)
			s.c <- err
			return
		}

		err = writeMessage(conn, req)
		if err != nil {
			s.closeConn(err)
			s.c <- err
			return
		}

		// signal message was written so method can return
		s.c <- nil

	reading:
		for {
			msg, err = readMessage(conn)
			if err != nil {
				tx <- StreamCompleted{Err: err}
				break
			}

			if res, ok := msg.(MessageResponse); ok {
				tx <- res
			} else {
				err := errors.New("protocol error: received request type")
				tx <- StreamCompleted{Err: err}
				break
			}

			switch msg.(type) {
			case StreamCompleted:
				writeMessage(conn, ClientAck{})
				break reading
			}
		}

		close(tx)
	}()

	select {
	case <-ctx.Done():
		s.closeConn(errCanceledCtx)
		<-s.c // drain channel
		return nil, errCanceledCtx
	case err := <-s.c:
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
}

func (s *Client) Close() error {
	// force panic if client calls Route/Health after Close
	defer close(s.c)

	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
