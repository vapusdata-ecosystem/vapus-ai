package pkgs

import (
	zerolog "github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
)

var DmLogger zerolog.Logger

func InitWAPLogger(debugMode bool) {
	DmLogger = dmlogger.GetDMLogger(debugMode, true, "")
}

var pkgLogger = dmlogger.GetSubDMLogger(DmLogger, IDEN, "pkgs")

func GetSubDMLogger(key, value string) zerolog.Logger {
	return dmlogger.GetSubDMLogger(DmLogger, key, value)
}
