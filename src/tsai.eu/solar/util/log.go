package util

import (
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
)

//------------------------------------------------------------------------------

// StartLogging initiate logging
func StartLogging() {
  // simple timeformat
  zerolog.TimeFieldFormat = ""

  // set log level
  config, err := GetConfiguration()
	if err != nil {
    LogError("main", "CORE", "unable to read the configuration")
	} else {
    switch config.CORE.LogLevel {
    case "panic":
      LogLevel(zerolog.PanicLevel)
    case "fatal":
      LogLevel(zerolog.FatalLevel)
    case "error":
      LogLevel(zerolog.ErrorLevel)
    case "warn":
      LogLevel(zerolog.WarnLevel)
    case "info":
      LogLevel(zerolog.InfoLevel)
    case "debug":
      LogLevel(zerolog.DebugLevel)
    default:
      LogLevel(zerolog.ErrorLevel)
    }
  }
}

//------------------------------------------------------------------------------

// LogLevel sets the LogLevel
//     panic (zerolog.PanicLevel, 5)
//     fatal (zerolog.FatalLevel, 4)
//     error (zerolog.ErrorLevel, 3)
//     warn  (zerolog.WarnLevel,  2)
//     info  (zerolog.InfoLevel,  1)
//     debug (zerolog.DebugLevel, 0)
func LogLevel(level zerolog.Level) {
  zerolog.SetGlobalLevel(level)
}

//------------------------------------------------------------------------------

// LogPanic logs panic information
func LogPanic(context string, module string, info string) {
  log.Panic().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------

// LogFatal logs fatal information
func LogFatal(context string, module string, info string) {
  log.Fatal().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------

// LogError logs error information
func LogError(context string, module string, info string) {
  log.Error().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------

// LogWarn logs warning information
func LogWarn(context string, module string, info string) {
  log.Warn().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------

// LogInfo logs info information
func LogInfo(context string, module string, info string) {
  log.Info().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------

// LogDebug logs debug information
func LogDebug(context string, module string, info string) {
  log.Debug().
    Str("Context", context).
    Str("Module", module).
    Msg(info)
}

//------------------------------------------------------------------------------
