package types

import (
	"time"
)

type DirectMessage struct {
	Id               int64  `json:"id"`
	Text             string `json:"text"`
	SenderId         int64  `json:"sender_id"`
	SenderScreenName string `json:"sender_screen_name"`
	createdAt        string `json:"created_at"`
}

func (dm *DirectMessage) CreatedAt() time.Time {
	at, _ := time.Parse(dm.createdAt, time.RubyDate)
	return at
}
