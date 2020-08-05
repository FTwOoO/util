package logging

import (
	"fmt"
	"github.com/rexue2019/util/errorkit"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

const (
	// Very verbose messages for debugging specific issues
	LevelDebug = "debug"
	// Default log level, informational
	LevelInfo = "info"
	// Warnings are messages about possible issues
	LevelWarn = "warn"
	// Errors are messages about things we know are problems
	LevelError = "error"
)

var Log Logger

const (
	KeyUserId       = "userId"
	KeyRoomId       = "roomId"
	KeyConnId       = "connId"
	KeyEvent        = "event"
	KeyMsg          = "msg"
	KeyScope        = "scope"
	KeyService      = "service"
	KeyScopeHTTP    = "http"
	KeyScopeMongoDb = "mongodb"
	KeyScopeGrpc    = "grpc"
)

func init() {
	Log = NewJsonLogger("info")
}

func SetLogSetting(level string, isJson bool) {
	if isJson {
		Log = NewJsonLogger(level)
	} else {
		Log = NewConsoleLogger(level)
	}
}

type Logger interface {
	GetLogLevel() string
	SetLogLevel(l string) error
	WithFields(fields map[string]interface{}) Logger //创建一个新logger，附带一些预设字段
	FatalError(error)
	LogError(error)
	Debugw(args ...interface{})
	Infow(args ...interface{})
	Warnw(args ...interface{})
	Errorw(args ...interface{})
	Fatalw(args ...interface{})
}

type loggerImp struct {
	zerolog.Logger
	logLevel zerolog.Level
}

func NewJsonLogger(loglevel string) Logger {
	var Logger = zerolog.New(os.Stderr)

	logLevelzero, err := zerolog.ParseLevel(loglevel)
	if err != nil {
		panic(err)
	}
	return &loggerImp{Logger: Logger, logLevel: logLevelzero}
}

func NewConsoleLogger(loglevel string) Logger {
	var Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	logLevelzero, err := zerolog.ParseLevel(loglevel)
	if err != nil {
		panic(err)
	}
	return &loggerImp{Logger: Logger, logLevel: logLevelzero}
}

func (this *loggerImp) WithFields(fields map[string]interface{}) Logger {
	newObj := &loggerImp{
		Logger:   this.Logger.With().Timestamp().Fields(fields).Logger(),
		logLevel: this.logLevel,
	}
	return newObj
}

func (this *loggerImp) GetLogLevel() string {
	return this.logLevel.String()
}

func (this *loggerImp) SetLogLevel(loglevel string) error {
	logLevelzero, err := zerolog.ParseLevel(loglevel)
	if err != nil {
		return err
	}
	this.Logger.WithLevel(logLevelzero)
	this.logLevel = logLevelzero
	return nil
}

func (this *loggerImp) FatalError(e error) {
	this.LogError(e)
	os.Exit(1)
}

func (this *loggerImp) LogError(err error) {
	e, ok := err.(errorkit.Error)
	if !ok {
		e = errorkit.WrapError(err)
	}

	zeroLogLevel, err := zerolog.ParseLevel(string(e.GetLogLevel()))
	if err != nil {
		this.Logger.Error().Msgf("invalid error:%v", e)
		return
	}

	if this.logLevel > zeroLogLevel {
		return
	}

	logEvent := this.Logger.WithLevel(zeroLogLevel)
	if e.GetMessage() != "" {
		logEvent = logEvent.Str("msg", e.GetMessage())
	}

	if len(e.GetParams()) > 0 {
		for k, v := range e.GetParams() {
			switch v.(type) {
			case string:
				logEvent = logEvent.Str(k, v.(string))
			case int:
				logEvent = logEvent.Int(k, v.(int))
			}
		}
	}
	if len(e.GetOp()) > 0 {
		logEvent = logEvent.Str("op", strings.Join(e.GetOp(), ","))
	}

	if len(e.GetLine()) > 0 {
		logEvent = logEvent.Str("line", e.GetLine())
	}

	if e.GetScope() != "" {
		logEvent = logEvent.Str("category", string(e.GetScope()))
	}
	logEvent = logEvent.Int("code", e.GetCode())
	logEvent.Msg("")
}

func (this *loggerImp) logW(logEvent *zerolog.Event, args ...interface{}) {
	if len(args)%2 != 0 {
		logEvent.Msg("invalidArgumentNums")
		return
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		value := args[i+1]
		switch value.(type) {
		case string:
			logEvent = logEvent.Str(key, value.(string))
		case int:
			logEvent = logEvent.Int(key, value.(int))
		case int8:
			logEvent = logEvent.Int8(key, value.(int8))
		case int16:
			logEvent = logEvent.Int16(key, value.(int16))
		case int32:
			logEvent = logEvent.Int32(key, value.(int32))
		case int64:
			logEvent = logEvent.Int64(key, value.(int64))
		case uint8:
			logEvent = logEvent.Uint8(key, value.(uint8))
		case uint16:
			logEvent = logEvent.Uint16(key, value.(uint16))
		case uint32:
			logEvent = logEvent.Uint32(key, value.(uint32))
		case uint64:
			logEvent = logEvent.Uint64(key, value.(uint64))
		default:
			logEvent = logEvent.Str(key, fmt.Sprintf("%v", value))
		}
	}

	logEvent.Msg("")
}

func (this *loggerImp) Debugw(args ...interface{}) {
	this.logW(this.Logger.Debug(), args...)
}

func (this *loggerImp) Infow(args ...interface{}) {
	this.logW(this.Logger.Info(), args...)
}

func (this *loggerImp) Warnw(args ...interface{}) {
	this.logW(this.Logger.Warn(), args...)
}

func (this *loggerImp) Errorw(args ...interface{}) {
	this.logW(this.Logger.Error(), args...)
}

func (this *loggerImp) Fatalw(args ...interface{}) {
	this.logW(this.Logger.Fatal(), args...)
}
