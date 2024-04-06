package sl

import "log/slog"

// Err adds an error field to slog record
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}