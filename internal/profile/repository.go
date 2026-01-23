package profile

import (
	"context"
	"fmt"

	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/uptrace/bun"
)

type RegistrationPayload struct {
	// Add fields as needed
}

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateProfile(ctx context.Context, profile *UserProfile, tx *bun.Tx) error {
	fmt.Println("CreateProfile repo is running")
	runner := utils.GetQueryRunner(tx, r.db)

	_, err := runner.NewInsert().Model(profile).Exec(ctx)
	fmt.Println("CreateProfile success")
	return err
}

func (r *Repository) GetProfileByUserId(ctx context.Context, userId string) (*UserProfile, error) {
	profile := new(UserProfile)
	err := r.db.NewSelect().Model(profile).Where("user_id = ?", userId).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
