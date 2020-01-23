package sonic

import (
	"encoding/json"
	"fmt"
	"io"
)

const lengthPrefixHeaderLen = 4

// TODO implement UnmarshalJSON MarshalJSON for protoMessage
// so we can use a much more efficient type as messageType
type messageType string

const (
	messageTypeAuth      messageType = "H"
	messageTypeCancel                = "C"
	messageTypeStarted               = "S"
	messageTypeQuery                 = "Q"
	messageTypeSchema                = "T"
	messageTypeProgress              = "P"
	messageTypeOutput                = "O"
	messageTypeAck                   = "A"
	messageTypeCompleted             = "D"
)

type protoMessage struct {
	MessageType messageType     `json:"e"`
	Variation   string          `json:"v,omitempty"`
	Payload     json.RawMessage `json:"p,omitempty"`
}

func decodeLen(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

func encodeLen(b []byte, length int32) {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(length >> 24 & 0x00FF)
	b[1] = byte(length >> 16 & 0x00FF)
	b[2] = byte(length >> 8 & 0x00FF)
	b[3] = byte(length)
}

func doWriteBuffer(writer io.Writer, msgBytes []byte) (err error) {
	buf := make([]byte, len(msgBytes)+lengthPrefixHeaderLen)
	encodeLen(buf[:lengthPrefixHeaderLen], int32(len(msgBytes)))
	copy(buf[lengthPrefixHeaderLen:], msgBytes)
	_, err = writer.Write(buf)
	return
}

func doReadBuffer(reader io.Reader) (buf []byte, err error) {
	var lengthPrefixBuf [lengthPrefixHeaderLen]byte

	_, err = io.ReadFull(reader, lengthPrefixBuf[:])
	if err != nil {
		return
	}

	length := decodeLen(lengthPrefixBuf[:])
	buf = make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		buf = nil
		return
	}
	return
}

func (m protoMessage) encode() []byte {
	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return b
}

func (m *protoMessage) decode(b []byte) (msg Message, err error) {
	err = json.Unmarshal(b, m)
	if err != nil {
		return
	}
	switch m.MessageType {
	case messageTypeAuth:
		msg, err = clientAuthFromProto(*m)
	case messageTypeCancel:
		msg, err = clientCancelFromProto(*m)
	case messageTypeStarted:
		msg, err = streamStartedFromProto(*m)
	case messageTypeQuery:
		msg, err = clientQueryFromProto(*m)
	case messageTypeSchema:
		msg, err = streamSchemaFromProto(*m)
	case messageTypeProgress:
		msg, err = streamProgressFromProto(*m)
	case messageTypeOutput:
		msg, err = streamOutputFromProto(*m)
	case messageTypeAck:
		msg, err = clientAckFromProto(*m)
	case messageTypeCompleted:
		msg, err = streamCompletedFromProto(*m)
	default:
		err = fmt.Errorf("unknown message type: %s", m.MessageType)
	}
	return
}
