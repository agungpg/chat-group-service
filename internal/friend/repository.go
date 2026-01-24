package friend

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateFriendRequest(ctx context.Context, fr FriendRequest) error {
	_, err := r.db.NewInsert().
		Model(&fr).
		Exec(ctx)

	return err
}

func (r *Repository) GetIncomingRequestList(ctx context.Context, userId string, limit, offset int) ([]FriendRequest, int, error) {
	frList := make([]FriendRequest, 0, limit)

	total, err := r.db.NewSelect().
		Model(&frList).
		Where("addressee_id = ?", userId).
		Where("status = ?", FriendRequestStatusPending).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	return frList, total, err
}

func (r *Repository) GetOutgoingRequestsList(ctx context.Context, userId string, limit, offset int) ([]FriendRequest, int, error) {
	frList := make([]FriendRequest, 0, limit)

	total, err := r.db.NewSelect().
		Model(&frList).
		Where("requester_id = ?", userId).
		Where("status = ?", FriendRequestStatusPending).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	return frList, total, err
}

func (r *Repository) AcceptRequest(ctx context.Context, id string) error {
	_, err := r.db.NewUpdate().
		Model(&FriendRequest{}).
		Set("updated_at = ?", time.Now()).
		Set("status = accepted").
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (r *Repository) DeclineRequest(ctx context.Context, id string) error {
	_, err := r.db.NewUpdate().
		Model(&FriendRequest{}).
		Set("updated_at = ?", time.Now()).
		Set("status = declined").
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (r *Repository) GetUserSummaries(ctx context.Context, userIDs []string) (map[string]userSummary, error) {
	if len(userIDs) == 0 {
		return map[string]userSummary{}, nil
	}

	rows := make([]userSummary, 0, len(userIDs))

	err := r.db.NewSelect().
		TableExpr("users AS u").
		Column("u.id", "u.username").
		ColumnExpr("COALESCE(up.display_name, u.username) AS display_name").
		ColumnExpr("up.avatar_url").
		Join("LEFT JOIN user_profiles AS up ON up.user_id = u.id").
		Where("u.id IN (?)", bun.In(userIDs)).
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}

	res := make(map[string]userSummary, len(rows))
	for _, row := range rows {
		res[row.ID] = row
	}

	return res, nil
}
