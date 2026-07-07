package user

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/Ayocodes24/GO-Eats/pkg/auth"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
	"github.com/golang-jwt/jwt/v5"
)

func (usrSrv *UsrService) Login(_ context.Context, userID int64, name string) (string, error) {
	claims := auth.UserClaims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "GO-Eats",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (usrSrv *UsrService) UserExist(ctx context.Context, email string, recordRequired bool) (bool, int64, string, error) {
	count, err := usrSrv.db.Count(ctx, "users", "COUNT(*)", "email", email)
	if err != nil {
		slog.Info("UserService.UserExist::Error %v", "error", err)
		return false, 0, "", err
	}
	if count == 0 {
		return false, 0, "", nil
	}

	if recordRequired {
		var accountInfo user.User
		err = usrSrv.db.Select(ctx, &accountInfo, "email", email)
		if err != nil {
			return false, 0, "", err
		}
		return true, accountInfo.ID, accountInfo.Name, nil
	}

	return true, 0, "", nil
}

func (usrSrv *UsrService) ValidatePassword(ctx context.Context, userInput *user.LoginUser) (bool, error) {
	var userAccount user.User
	err := usrSrv.db.Select(ctx, &userAccount, "email", userInput.Email)
	if err != nil {
		slog.Info("UserService.ValidatePassword::Error %v", "error", err)
		return false, err
	}

	err = userInput.CheckPassword(userAccount.Password)
	if err != nil {
		return false, errors.New("invalid password")
	}
	return true, nil
}
