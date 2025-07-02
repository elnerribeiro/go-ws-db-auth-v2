package utils

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type LoggerContext context.Context

var logContext *LoggerContext
var logger *zerolog.Logger

// Message returns a Json message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond encodes the response
func Respond(logContext *LoggerContext, w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger := zerolog.Ctx(*logContext)
		logger.Error().Err(err).Msg("Error while encoding response")
		return
	}
}

func GetLoggerAndContext() (*zerolog.Logger, *LoggerContext) {
	return logger, logContext
}

func init() {
	initLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Caller().
		Logger()
	initLogContext := context.Background()
	initLogContext = initLogger.WithContext(initLogContext)
	logger = &initLogger
	logContext = (*LoggerContext)(&initLogContext)
}
