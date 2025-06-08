package repository

import (
	"context"
	"fmt"
	"gop_shlyop/internal/types"
	"strings"
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

// retrieves a list of reviews, optionally filtered by user_id or item_id.
func (r *ReviewsRepository) GetReviews(ctx context.Context, userID, itemID int64) ([]types.Review, error) {
	baseQuery := `
		SELECT id, user_id, item_id, text, sentiment, created_at
		FROM reviews
	`
	var args []interface{}
	var conditions []string

	if userID > 0 {
		args = append(args, userID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if itemID > 0 {
		args = append(args, itemID)
		conditions = append(conditions, fmt.Sprintf("item_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []types.Review
	for rows.Next() {
		var review types.Review
		if err := rows.Scan(&review.ID, &review.UserID, &review.ItemID, &review.Text, &review.Sentiment, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan review row: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading review rows: %w", err)
	}

	if reviews == nil {
		reviews = make([]types.Review, 0)
	}

	return reviews, nil
}
