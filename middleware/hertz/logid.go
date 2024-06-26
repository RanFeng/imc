package hertz

import (
	"context"

	"github.com/RanFeng/logid"
	"github.com/cloudwego/hertz/pkg/app"
)

const (
	LogIDKey = "K_LOG_ID"
)

// InjectLogID 为hertz的上下文注入logid
func InjectLogID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logID, _ := ctx.Value(LogIDKey).(string)
		if len(logID) == 0 {
			ctx = context.WithValue(ctx, LogIDKey, logid.GenLogID())
		}
		c.Next(ctx)
	}
}
