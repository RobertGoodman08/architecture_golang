package context

import (
	"context"
	"time"
)

type Value interface {
	Value(key any) any
	WithValue(key, value any)

	ID() string
}

func (l *local) ID() string {
	id, _ := l.id()
	return id
}

func (l *local) id() (string, bool) {
	value := l.Value(keyRequestID)
	id, ok := value.(string)
	return id, ok
}

func (l *local) Value(key any) any {
	return l.base.Value(key)
}

func (l *local) WithValue(key, value any) {
	if key == keyRequestID {
		return // ignore
	}

	l.withValue(key, value)
}

func (l *local) withValue(key, value any) {
	l.base = context.WithValue(l.base, key, value)
}

func (l *local) WithTimeout(timeout time.Duration) {
	l.base, l.cancelFunc = context.WithTimeout(l.base, timeout)
}

func (l local) CopyWithTimeout(timeout time.Duration) Context {
	l.base, l.cancelFunc = context.WithTimeout(l.base, timeout)
	return &l
}

func (l *local) WithDeadline(d time.Time) {
	l.base, l.cancelFunc = context.WithDeadline(l.base, d)
}

func (l local) CopyWithDeadline(d time.Time) Context {
	l.base, l.cancelFunc = context.WithDeadline(l.base, d)
	return &l
}
