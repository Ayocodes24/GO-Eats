package utils

import "time"

type Timestamp struct {
	CreatedAt time.Time `bun:",null zero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",null zero,notnull,default:current_timestamp"`
}
