package user

import (
	"context"
	"fmt"

	"github.com/thearyanahmed/endframe/pkg/entity"
)

type UserRepositoryInterface interface {
	FindById(ctx context.Context, id int) (*entity.User, error)
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindById(ctx context.Context, id int) (*entity.User, error) {
	return &entity.User{}, fmt.Errorf("not yet implemented")
}
