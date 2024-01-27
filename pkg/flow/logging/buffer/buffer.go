package buffer

import (
	"context"
	"log/slog"
	"sync"

	"github.com/go-kit/log"
	gokitlevel "github.com/go-kit/log/level"
	"github.com/grafana/agent/pkg/flow/logging"
)

const (
	levelKey = "level"
)

type Buffer struct {
	mut          sync.Mutex
	hasLogFormat bool
	buffer       [][]interface{}
}

var Logger = &Buffer{
	hasLogFormat: false,
	buffer:       [][]interface{}{},
}

func (b *Buffer) LogInfo(logger log.Logger, keyvals ...interface{}) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.hasLogFormat {
		toLevel(logger, gokitlevel.ErrorValue(), slog.LevelInfo).Log(keyvals...)
		return
	}
	b.buffer = append(b.buffer, keyvals)
}

func (b *Buffer) LogError(logger log.Logger, keyvals ...interface{}) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.hasLogFormat {
		toLevel(logger, gokitlevel.ErrorValue(), slog.LevelError).Log(keyvals...)
		return
	}
	b.buffer = append(b.buffer, keyvals)
}

func (b *Buffer) LogWarn(logger log.Logger, keyvals ...interface{}) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.hasLogFormat {
		toLevel(logger, gokitlevel.ErrorValue(), slog.LevelWarn).Log(keyvals...)
		return
	}
	b.buffer = append(b.buffer, keyvals)
}

func (b *Buffer) LogDebug(logger log.Logger, keyvals ...interface{}) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.hasLogFormat {
		toLevel(logger, gokitlevel.ErrorValue(), slog.LevelDebug).Log(keyvals...)
		return
	}
	b.buffer = append(b.buffer, keyvals)
}

func (b *Buffer) LogBuffered(logger log.Logger) {
	b.mut.Lock()
	if !b.hasLogFormat {
		return
	}
	for _, buffered := range b.buffer {
		logger.Log(buffered...)
	}
	b.mut.Unlock()
}

func (b *Buffer) HasLogFormatNow() {
	b.mut.Lock()
	b.hasLogFormat = true
	b.mut.Unlock()
}

func toLevel(logger log.Logger, level gokitlevel.Value, slogLevel slog.Level) log.Logger {
	switch l := logger.(type) {
	case logging.EnabledAware:
		if !l.Enabled(context.Background(), slogLevel) {
			return disabledLogger
		}
	}
	return log.WithPrefix(logger, levelKey, level)
}

var disabledLogger = &noopLogger{}

type noopLogger struct{}

func (d *noopLogger) Log(_ ...interface{}) error {
	return nil
}
