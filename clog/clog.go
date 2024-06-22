// Package clog provides a slog Handler that outputs colored logs.
package clog

import (
	"context"
	"hash/fnv"
	"io"
	"log/slog"
)

const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	white   = "\033[97m"
)

var allColors = []string{
	red,
	green,
	yellow,
	blue,
	magenta,
	cyan,
	gray,
	white,
}

// Handler embeds a text handler (maybe just needs to embed a handler)
// so that it implements a [slog.Handler], but overrides methods so that the output
// has color.
type Handler struct {
	w io.Writer
	*slog.TextHandler

	svc   string
	color string
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func NewHandler(w io.Writer, svcName string, opts *slog.HandlerOptions) *Handler {
	// Determine what color to make this service's logs.
	// It's determinant so the services will always have the same colors.
	n := int(hash(svcName))

	return &Handler{
		w:           w,
		TextHandler: slog.NewTextHandler(w, opts),
		svc:         svcName,
		color:       allColors[n%len(allColors)],
	}
}

// Handle wraps the output of the log with an ANSI color code.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("service", h.svc))

	h.w.Write([]byte(h.color))
	err := h.TextHandler.Handle(ctx, r)
	h.w.Write([]byte(reset))

	return err
}
