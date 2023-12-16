package presenter

import (
	"github.com/thearyanahmed/endframe/pkg/entity"
)

type UserResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"Name"`
}

func ToUserResponse(userEntity entity.User) *UserResponse {
	return &UserResponse{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}
}
