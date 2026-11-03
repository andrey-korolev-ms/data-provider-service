package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	greetv1 "github.com/andrey-korolev-ms/data-provider-service/gen/greet/v1"
	"github.com/andrey-korolev-ms/data-provider-service/internal/delivery"
	"github.com/andrey-korolev-ms/data-provider-service/internal/domain"
	"github.com/andrey-korolev-ms/data-provider-service/internal/interceptor"
	"github.com/andrey-korolev-ms/data-provider-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Инициализация Logger (JSON формат удобно для парсинга в ELK/Loki)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	// сборка слоев
	greeterDomain := domain.NewSimpleGreeter()
	greeterUsecase := usecase.NewGreeterUsecase(greeterDomain)
	grpcAdapter := delivery.NewGrpcServer(greeterUsecase)

	// создаем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	// 4. Создание gRPC сервера С ИНТЕРСЕПТОРАМИ
	// ВАЖНО: Порядок имеет значение!
	// Recovery должен быть самым внешним (первым), чтобы поймать панику от любого другого слоя.
	// Logging может быть вторым, чтобы записать успешные/неуспешные попытки.

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnaryServerInterceptor(logger),
			interceptor.LoggingUnaryServerInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			interceptor.RecoveryStreamServerInterceptor(logger),
			interceptor.LoggingStreamServerInterceptor(logger),
		),
	}
	s := grpc.NewServer(opts...)

	// регистрируем сервис
	greetv1.RegisterGreeterServiceServer(s, grpcAdapter)
	// включаем reflection (чтобы можно было использовать grpcurl без .proto файла на клиенте)
	reflection.Register(s)

	log.Println("gRPC started on 50051 with interceptor")
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err)
	}
}
