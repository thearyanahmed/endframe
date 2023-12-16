package user

import (
	"context"
)

type userServiceInterface interface {
	Find(ctx context.Context, id int) string
}

type UserService struct {
	repository UserRepositoryInterface
}

func NewUserService(repository UserRepositoryInterface) *UserService {
	return &UserService{
		repository,
	}
}

func (svc *UserService) Find(ctx context.Context, id int) string {
	user, err := svc.repository.FindById(ctx, id)

	if err != nil {
		return ""
	}

	return user.Name
}
