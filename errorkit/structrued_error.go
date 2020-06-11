package errorkit

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorScope string

const (
	ErrorScopeFile    ErrorScope = "file"
	ErrorScopeSocket  ErrorScope = "socket"
	ErrorScopeMongoDb ErrorScope = "mongodb"
	ErrorScopeMySQL   ErrorScope = "mysql"
	ErrorScopeRedis   ErrorScope = "redis"
	ErrorScopeHttp    ErrorScope = "http"
)

type ErrorLogLevel string

const (
	ErrorLogLevelDebug ErrorLogLevel = "debug"
	ErrorLogLevelInfo  ErrorLogLevel = "info"
	ErrorLogLevelWarn  ErrorLogLevel = "warn"
	ErrorLogLevelError ErrorLogLevel = "error"
	ErrorLogLevelFatal ErrorLogLevel = "fatal"
)

/*type StructuredError interface {
	//层层添加错误调用栈
	AddOp(op string) StructuredError
	//层层添加参数
	AddParam(k string, v interface{}) StructuredError
}*/

type StructuredError struct {
	Op       []string               //方法调用链
	Line     string                 //最底层的调用行
	Scope    ErrorScope             //错误分类
	Code     int                    //唯一错误码，如果没有定义，则是-1
	Params   map[string]interface{} //自定义参数
	Err      error                  //底层错误
	Message  string                 //自定义错误时Err为nil，但需要添加一条文本作为错误详细信息（供后台查看）
	LogLevel ErrorLogLevel          // debug/info/warn/error这几个日志级别
}

func GetCallFrame(skip int) runtime.Frame {
	var pc [1]uintptr
	n := runtime.Callers(skip, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}

func WrapError(err error) Error {
	if er, ok := err.(Error); ok {
		return er
	}

	frame := GetCallFrame(3)
	return &StructuredError{
		Line:     fmt.Sprintf("%s:%d", frame.File, frame.Line),
		LogLevel: ErrorLogLevelError,
		Err:      err,
		Op:       []string{funcname(frame.Function)},
	}
}

func NewStructuredError() Error {
	frame := GetCallFrame(3)
	return &StructuredError{
		Line:     fmt.Sprintf("%s:%d", frame.File, frame.Line),
		LogLevel: ErrorLogLevelError,
		Op:       []string{funcname(frame.Function)},
	}
}

func (this *StructuredError) Error() string {
	return this.GetMessage()
}

func (this *StructuredError) SetLogLevel(level ErrorLogLevel) Error {

	this.LogLevel = level
	return this
}

func (this *StructuredError) SetError(err error) Error {
	this.Err = err
	return this
}

func (this *StructuredError) AddOp(op string) Error {
	this.Op = append(this.Op, op)
	return this
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func (this *StructuredError) AddFunctionAsOp() Error {
	frame := GetCallFrame(3)
	this.Op = append(this.Op, funcname(frame.Function))
	return this
}

func (this *StructuredError) SetCode(code int) Error {
	this.Code = code
	return this
}

func (this *StructuredError) SetMessage(msg string) Error {
	this.Message = msg
	return this
}

func (this *StructuredError) SetScope(c ErrorScope) Error {
	this.Scope = c
	return this
}

func (this *StructuredError) AddParam(k string, v interface{}) Error {
	if this.Params == nil {
		this.Params = make(map[string]interface{})
	}
	this.Params[k] = v
	return this
}

func (this *StructuredError) GetOp() []string {
	return this.Op
}

func (this *StructuredError) GetScope() ErrorScope {
	return this.Scope
}

func (this *StructuredError) GetCode() int {
	return this.Code
}

func (this *StructuredError) GetParams() map[string]interface{} {
	return this.Params
}

func (this *StructuredError) GetMessage() string {
	if this.Message != "" {
		return this.Message
	} else if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

func (this *StructuredError) GetLogLevel() ErrorLogLevel {
	return this.LogLevel
}

func (this *StructuredError) GetLine() string {
	return this.Line
}
