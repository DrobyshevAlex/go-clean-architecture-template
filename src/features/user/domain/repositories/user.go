package repositories

import (
	"context"
	"main/src/features/user/domain/entities"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int64) (*entities.User, error)
}
