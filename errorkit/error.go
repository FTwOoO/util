package errorkit

type Error interface {
	error
	SetLogLevel(level ErrorLogLevel) Error
	SetError(err error) Error
	AddOp(op string) Error
	AddFunctionAsOp() Error
	SetCode(code int) Error
	SetMessage(msg string) Error
	SetScope(c ErrorScope) Error
	SetEvent(e string) Error
	AddParam(k string, v interface{}) Error

	GetLine() string
	GetOp() []string
	GetScope() ErrorScope
	GetCode() int
	GetParams() map[string]interface{}
	GetMessage() string
	GetLogLevel() ErrorLogLevel
}

type HttpError interface {
	Error
	GetHttpCode() int
}
