package webapp

import "log/slog"

type Logger interface {
}

func init() {
	l := slog.Logger{}
	l.With()
}
