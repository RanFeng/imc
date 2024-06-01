package hertz

import (
	"context"
	"os"
	"time"

	"github.com/RanFeng/ilog"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "counter",
		Help: "A simple counter metric",
	}, []string{"psm", "method", "path", "status", "log_id"})

	latency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "latency",
		Help: "Duration of HTTP requests",
		Objectives: map[float64]float64{
			0.5:  0.1,  // 50% 分位数,最大 10% 误差
			0.9:  0.05, // 90% 分位数,最大 5% 误差
			0.99: 0.01, // 99% 分位数,最大 1% 误差
		},
	}, []string{"psm", "method", "path", "status", "log_id"})

	pushServer *push.Pusher
	psm        string
)

func init() {
	prometheus.MustRegister(counter)
	pushServer = push.New("http://localhost:18974", "common_metrics")
	psm = os.Getenv("PSM")
}

// CommonMetrics 为hertz的接口记录通用的指标
func CommonMetrics() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		logID, _ := ctx.Value(LogIDKey).(string)
		vList := []string{psm,
			string(c.Method()),
			string(c.Path()),
			string(rune(c.Response.Header.StatusCode())),
			logID}
		counter.WithLabelValues(vList...).Inc()
		latency.WithLabelValues(vList...).Observe(float64(time.Since(start).Milliseconds()))

		err := pushServer.
			Collector(counter).
			Collector(latency).
			Push()
		if err != nil {
			ilog.EventWarn(ctx, "push_common_metrics_fail", "v_list", vList, "err", err)
		}
	}
}
