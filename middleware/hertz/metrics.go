package hertz

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/RanFeng/ilog"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	keys    = []string{"psm", "method", "path", "status"}
	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "imc_middleware_hertz_common_metrics_counter",
		Help: "用于记录接口请求量",
	}, keys)

	latency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "imc_middleware_hertz_common_metrics_latency",
		Help: "用于记录接口请求时延",
		Objectives: map[float64]float64{
			0.5:  0.1,  // 50% 分位数,最大 10% 误差
			0.9:  0.05, // 90% 分位数,最大 5% 误差
			0.99: 0.01, // 99% 分位数,最大 1% 误差
		},
	}, keys)

	pushServer *push.Pusher
	psm        string
)

func init() {
	pushServer = push.New("http://localhost:18976", "common_metrics").
		Collector(counter).Collector(latency)
	psm = os.Getenv("PSM")
}

// CommonMetrics 为hertz的接口记录通用的指标
func CommonMetrics() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		vList := []string{psm,
			string(c.Method()),
			c.FullPath(),
			strconv.Itoa(c.Response.Header.StatusCode())}
		counter.WithLabelValues(vList...).Inc()
		latency.WithLabelValues(vList...).Observe(float64(time.Since(start).Milliseconds()))

		err := pushServer.Push()
		if err != nil {
			ilog.EventWarn(ctx, "push_common_metrics_fail", "v_list", vList, "err", err)
		}
	}
}
