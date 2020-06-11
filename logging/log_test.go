package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.livedev.shika2019.com/go/util/errorkit"
	"os"
	"testing"
)

func TestSetLogLevel(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("foo", "bar").Msg("")
}

func TestLogger2_LogError(t *testing.T) {
	logger := NewConsoleLogger("info")
	er := errorkit.NewStructuredError().
		AddParam("category", errorkit.ErrorScopeMongoDb).
		SetCode(111).
		AddParam("ticketsCount", 1).AddParam("userId", "112233").SetError(fmt.Errorf("file not open"))
	logger.LogError(er)
}

func TestLogger2_LogEvent(t *testing.T) {
	logger := NewJsonLogger("info")
	logger.Infow(
		KeyUserId, "111",
		KeyEvent, "help",
	)
}

func TestLoggerImp_WithFields(t *testing.T) {
	logger := NewJsonLogger("info").WithFields(map[string]interface{}{
		KeyScope:   "mongodb",
		KeyService: "im",
	})
	logger.Infow("key", "haha")
}
