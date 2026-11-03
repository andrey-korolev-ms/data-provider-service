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
func LoggingStreamServerInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		// вызываем обрабочик стрима
		err := handler(srv, ss)
		duration := time.Since(start)

		logAttrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", duration),
		}
		if err != nil {
			st, _ := status.FromError(err)
			logAttrs = append(logAttrs,
				slog.Int("code", int(st.Code())),
				slog.String("message", st.Message()),
			)
			logger.Error("gRPC stream failed", logAttrs...)
		} else {
			logger.Info("gRPC stream finished ", logAttrs...)
		}
		return err
	}
}
