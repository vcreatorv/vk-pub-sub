package utils

import (
	"context"
	"time"
)

type ctxTimeout string

var ctxTimeoutKey ctxTimeout = "timeout"

func GetTimeout(ctx context.Context) time.Duration {
	if timeout, ok := ctx.Value(ctxTimeoutKey).(time.Duration); ok {
		return timeout
	}
	return 5 * time.Second
}

func SetTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return context.WithValue(ctx, ctxTimeoutKey, timeout)
}
