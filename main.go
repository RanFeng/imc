package imc

import "github.com/RanFeng/imc/middleware/hertz"

// export
var (
	InjectLogID   = hertz.InjectLogID
	CommonMetrics = hertz.CommonMetrics
)
