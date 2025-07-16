package user

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `bun:",pk,autoincrement" json:"id"`
	Name          string `bun:",notnull" json:"name" validate:"name"`
	Email         string `bun:",unique,notnull" json:"email" validate:"email"`
	PasswordHash  string `bun:",notnull" json:"password_hash"`
}
