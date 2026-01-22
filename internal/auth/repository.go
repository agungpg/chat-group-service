package auth

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("username = ?", username).Where("is_deleted = FALSE").Scan(ctx)

	return user, err
}
