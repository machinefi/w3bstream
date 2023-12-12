package consts

import "log/slog"

//go:generate toolkit gen enum LogLevel
type LogLevel int32

const (
	LOG_LEVEL_UNKNOWN LogLevel = iota

	_ // FATAL
	LOG_LEVEL__ERROR
	LOG_LEVEL__WARN
	LOG_LEVEL__INFO
	LOG_LEVEL__DEBUG
	_ // TRACE
)

func (v LogLevel) Level() slog.Level {
	switch v {
	case LOG_LEVEL__ERROR:
		return slog.LevelError
	case LOG_LEVEL__WARN:
		return slog.LevelWarn
	case LOG_LEVEL__INFO:
		return slog.LevelInfo
	case LOG_LEVEL__DEBUG:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
