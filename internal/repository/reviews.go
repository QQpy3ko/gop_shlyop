package repository

import (
	"context"
	"fmt"
	"gop_shlyop/internal/types"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewsRepository struct {
	db *pgxpool.Pool
}

func NewReviewsRepository(db *pgxpool.Pool) *ReviewsRepository {
	return &ReviewsRepository{db: db}
}

// saves a new review to the database.
func (r *ReviewsRepository) CreateReview(ctx context.Context, req types.AddReviewRequest, sentiment string) (types.Review, error) {
	query := `
		INSERT INTO reviews (user_id, item_id, text, sentiment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var reviewID int64
	var createdAt time.Time

	err := r.db.QueryRow(ctx, query, req.UserID, req.ItemID, req.Text, sentiment).Scan(&reviewID, &createdAt)
	if err != nil {
		return types.Review{}, fmt.Errorf("failed to create review: %w", err)
	}

	return types.Review{
		ID:        reviewID,
		UserID:    req.UserID,
		ItemID:    req.ItemID,
		Text:      req.Text,
		Sentiment: sentiment,
		CreatedAt: createdAt,
	}, nil
}
