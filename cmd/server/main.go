package main

import (
	"log"
	"net"

	greetv1 "github.com/andrey-korolev-ms/data-provider-service/gen/greet/v1"
	"github.com/andrey-korolev-ms/data-provider-service/internal/delivery"
	"github.com/andrey-korolev-ms/data-provider-service/internal/domain"
	"github.com/andrey-korolev-ms/data-provider-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// инициализация слоев
	greeterDomain := domain.NewSimpleGreeter()
	greeterUsecase := usecase.NewGreeterUsecase(greeterDomain)
	grpcAdapter := delivery.NewGrpcServer(greeterUsecase)

	// создаем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// регистрируем сервис
	greetv1.RegisterGreeterServiceServer(s, grpcAdapter)
	// включаем reflection (чтобы можно было использовать grpcurl без .proto файла на клиенте)
	reflection.Register(s)

	log.Println("gRPC started on 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
