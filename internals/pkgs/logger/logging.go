package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type DMLogger struct {
	zerolog.Logger
}

var DmLogg *DMLogger

var CoreLogger = GetDMLogger(false, true, "")

// New creates a new logger with the given options.
func NewZeroLogger(debugLogFlag bool, consoleWritter bool, withCaller bool, separator string, exclude ...string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLogFlag {

		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + "::" + strconv.Itoa(line)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if !consoleWritter {
		l := zerolog.Logger{}
		if withCaller {
			return l.With().Caller().Logger()
		} else {
			return l
		}
	}

	if exclude == nil {
		exclude = []string{}
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, PartsExclude: exclude}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(" %s %v ", i, separator)
	}

	// Format level: fatal, error, debug, info, warn
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s : ", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s %v ", i, separator)
	}

	output.FormatCaller = func(i interface{}) string {

		return strings.ToUpper(fmt.Sprintf("%-6s%v", i, separator))
	}

	// format error
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf(" %s : ", i)
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s %v ", i, separator)

	}
	// format timestamp
	output.TimeFormat = time.RFC3339
	if withCaller {
		return zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		return zerolog.New(output).With().Timestamp().Logger()
	}
}

func GetDMLogger(debugLogFlag bool, consoleWritter bool, exclude ...string) zerolog.Logger {
	return NewZeroLogger(debugLogFlag, consoleWritter, true, "|", exclude...)
}

func GetSubDMLogger(logger zerolog.Logger, key, value string) zerolog.Logger {
	return logger.With().Str(key, value).Logger()
}

/*
TBD
*/
func (dml *DMLogger) LogWL(msg string, vals map[string]string) {
	dml.Logger.Info().Fields(vals).Msg(msg)
}

func (dml *DMLogger) LogErrorWithFields(err error, msg string, vals map[string]string) {
	dml.Logger.Error().Err(err).Fields(vals).Msg(msg)
}

func (dml *DMLogger) LogErrorWithStacktrace(err error, msg string) {
	dml.Logger.Error().Stack().Err(err).Msg(msg)
}

func (dml *DMLogger) LogError(err error, msg string) {
	dml.Logger.Error().Err(err).Msg(msg)
}

func (dml *DMLogger) LogFatal(err error, msg string) {
	dml.Logger.Fatal().Err(err).Msg(msg)
}

func (dml *DMLogger) LogFileInfo(msg string) {
	dml.Logger = log.With().Caller().Logger()
	dml.Logger.Info().Msg(msg)
}

func (dml *DMLogger) LogInfo(msg string) {
	dml.Logger.Info().Msg(msg)
}

func (dml *DMLogger) LogDebug(msg string) {
	dml.Logger.Debug().Msg(msg)
}

func (dml *DMLogger) LogFInfo(msg string, data ...interface{}) {
	dml.Logger.Info().Msgf(msg, data...)
}

func (dml *DMLogger) LogFDebug(msg string, data ...interface{}) {
	dml.Logger.Debug().Msgf(msg, data...)
}

func (dml *DMLogger) LogErrorWithConext(ctx context.Context, err error, msg string) {
	dml.Logger.Error().Ctx(ctx).Err(err).Msg(msg)
}

func (dml *DMLogger) LogWInfoFields(msg string, fields map[string]any) {
	dml.Logger.Info().Fields(fields).Msg(msg)
}

func (dml *DMLogger) LogWErrFields(msg string, fields map[string]any) {
	dml.Logger.Error().Fields(fields).Msg(msg)
}

func (dml *DMLogger) LogWDebugFields(msg string, fields map[string]any) {
	dml.Logger.Debug().Fields(fields).Msg(msg)
}
