package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)
		duration := time.Since(start)
		// вызываем следующий обработчик в цепочке (или сам метод сервиса)
		logAttrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("deration", duration),
		}
		if err != nil {
			st, _ := status.FromError(err)
			logAttrs = append(logAttrs,
				slog.Int("code", int(st.Code())),
				slog.String("message", st.Message()),
			)
			logger.Error("gRPC requst failed", logAttrs...)
		} else {
			logger.Info("gRPC requst succeeded", logAttrs...)
		}
		return resp, err
	}
}
