package types

import (
	"database/sql"
	"time"
)

type Review struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ItemID    int64     `json:"item_id"`
	Text      string    `json:"text"`
	Sentiment string    `json:"sentiment"`
	CreatedAt time.Time `json:"created_at"`
}

type AddReviewRequest struct {
	UserID int64  `json:"user_id"`
	ItemID int64  `json:"item_id"`
	Text   string `json:"text"`
}
