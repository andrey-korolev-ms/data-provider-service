package usecase

import (
	"context"
	"errors"

	"github.com/andrey-korolev-ms/data-provider-service/internal/domain"
)

type GreeterUsecase interface {
	HandleSayHello(ctx context.Context, name string) (string, error)
	HandleManyHellos(ctx context.Context, name string, count int) ([]string, error)
}

type greeterUsecase struct {
	repo domain.Greeter
}

func NewGreeterUsecase() *greeterUsecase {
	return &greeterUsecase{}
}
func (u *greeterUsecase) HandleSayHello(ctx context.Context, name string) (string, error) {
	//какая-то допустим валидация
	if len(name) > 50 {
		return "", errors.New("name too long")
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return u.repo.SayHello(name), nil
	}
}
func (u *greeterUsecase) HandleManyHellos(ctx context.Context, name string, count int) ([]string, error) {
	if count <= 0 || count > 100 {
		return nil, errors.New("invalid number of replies")
	}
	return u.repo.GenerateHellos(name, count), nil
}
