package delivery

import (
	"context"
	"log"

	greetv1 "github.com/andrey-korolev-ms/data-provider-service/gen/greet/v1"
	"github.com/andrey-korolev-ms/data-provider-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	greetv1.UnimplementedGreeterServiceServer // заглушка на будущее
	uc                                        usecase.GreeterUsecase
}

func NewGrpcServer(uc usecase.GreeterUsecase) *GrpcServer {
	return &GrpcServer{uc: uc}
}
func (s *GrpcServer) SayHello(ctx context.Context, req *greetv1.HelloRequest) (*greetv1.HelloReply, error) {
	log.Printf("Получен одиночный запрос для name: %s", req.GetName())
	msg, err := s.uc.HandleSayHello(ctx, req.GetName())
	if err != nil {
		//преобразуем Go-error в gRPC status
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &greetv1.HelloReply{Message: msg}, nil
}
func (s *GrpcServer) ManyHellos(req *greetv1.ManyHellosRequest, stream greetv1.GreeterService_ManyHellosServer) error {
	log.Printf("Запуск стрима для name: %s, count %d", req.GetName(), req.NumReplies)
	// получаем данные из бизнес-логики
	messages, err := s.uc.HandleManyHellos(stream.Context(), req.GetName(), int(req.NumReplies))
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	for _, msg := range messages {
		reply := &greetv1.HelloReply{Message: msg}
		if err := stream.Send(reply); err != nil {
			return status.Errorf(codes.Canceled, "client canceled stream: %v", err)
		}
	}
	return nil
}
