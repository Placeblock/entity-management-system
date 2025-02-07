package models

import "time"

type Token struct {
	EntityID  int64
	Pin       string
	CreatedAt time.Time
}
