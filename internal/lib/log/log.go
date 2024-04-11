package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func Init(level slog.Level) *slog.Logger {
	opts := HandlerOptions{
		SlogOptions: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewHandler(os.Stdout)

	return slog.New(handler)
}

type HandlerOptions struct {
	SlogOptions *slog.HandlerOptions
}

type Handler struct {
	slog.Handler
	attrs []slog.Attr
}

func (opts HandlerOptions) NewHandler(w io.Writer) *Handler {
	h := &Handler{
		Handler: slog.NewJSONHandler(w, opts.SlogOptions),
	}

	return h
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	line := CompleteLogMessage(r)
	filename := fmt.Sprintf(".logs/api-logs-%s.logs", r.Time.Format("02.01.2006"))

	fmt.Printf(line)
	err := ToFile(filename, CompleteLogMessage(r))
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		Handler: h.Handler,
		attrs:   attrs,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Handler: h.Handler.WithGroup(name),
	}
}

func CompleteLogMessage(r slog.Record) string {
	line := fmt.Sprintf(`{"time":"%s","level":"%s","message":"%s"`, r.Time.Format("02-01-2006T15:04:05Z07:00"), r.Level.String(), r.Message)
	r.Attrs(func(a slog.Attr) bool {
		line += fmt.Sprintf(`,"%s":%v`, a.Key, a.Value.Any())
		return true
	})
	line += "}\n"

	return line
}

func ToFile(filename, message string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(message))
	if err != nil {
		return err
	}
	return f.Close()
}
