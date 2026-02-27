package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

func AddTraceID(ctx context.Context) context.Context {
	// 生成唯一的trace_id使用标准库
	if _, ok := ctx.Value(TraceIDKey{}).(string); ok {
		return ctx
	}
	randBytes := make([]byte, 3)
	rand.Read(randBytes)
	randomStr := hex.EncodeToString(randBytes)
	traceId := "trace-" + strconv.FormatInt(time.Now().UnixNano(), 10)[8:] + "-" + strings.ToLower(randomStr)
	ctx = context.WithValue(ctx, TraceIDKey{}, traceId)
	return ctx
}
