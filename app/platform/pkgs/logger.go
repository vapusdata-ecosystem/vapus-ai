package pkgs

import (
	zerolog "github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

var DmLogger zerolog.Logger
var pkgLogger zerolog.Logger

func InitWAPLogger(debugMode bool) {
	DmLogger = dmlogger.GetDMLogger(debugMode, true, "")
	pkgLogger = dmlogger.GetSubDMLogger(DmLogger, IDEN, "pkgs")
}

func GetSubDMLogger(key, value string) zerolog.Logger {
	return dmlogger.GetSubDMLogger(DmLogger, key, value)
}

func getSubPkgLogger(key, value string) zerolog.Logger {
	return dmlogger.GetSubDMLogger(pkgLogger, key, value)
}
