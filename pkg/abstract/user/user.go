package user

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
)

type User interface {
	Add(ctx context.Context, user *user.User) (bool, error)
	Delete(ctx context.Context, userId int64) (bool, error)
	Login(ctx context.Context, userID int64) (string, error)
	UserExist(ctx context.Context, email string, recordRequired bool) (bool, int64, error)
}
