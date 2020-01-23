package sonic

import (
	logd "github.com/ernestrc/logd-go/logging"
)

const (
	keyThread        string = logd.KeyThread
	keyClass                = logd.KeyClass
	keyLevel                = logd.KeyLevel
	keyTime                 = logd.KeyTime
	keyDate                 = logd.KeyDate
	keyTimestamp            = logd.KeyTimestamp
	keyMessage              = logd.KeyMessage
	keyCallType             = "callType"
	keyStep                 = "step"
	keyTraceID              = "traceID"
	keyTraceIDSnake         = "trace_id"
	valueStepSuccess        = "success"
	valueStepFailure        = "failure"
	valueStepAttempt        = "attempt"
	keyError                = "error"
)
