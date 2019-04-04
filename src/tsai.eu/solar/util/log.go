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
    LogLevel(config.CORE.LogLevel)
  }
}

//------------------------------------------------------------------------------

// LogLevel sets the LogLevel
func LogLevel(level string) {
  switch level {
  case "panic":
    zerolog.SetGlobalLevel(zerolog.PanicLevel)
  case "fatal":
    zerolog.SetGlobalLevel(zerolog.FatalLevel)
  case "error":
    zerolog.SetGlobalLevel(zerolog.ErrorLevel)
  case "warn":
    zerolog.SetGlobalLevel(zerolog.WarnLevel)
  case "info":
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
  case "debug":
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
  default:
    zerolog.SetGlobalLevel(zerolog.ErrorLevel)
  }
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
