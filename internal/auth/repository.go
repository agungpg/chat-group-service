package auth

import (
	"context"
	"time"

	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User, tx *bun.Tx) error {

	runner := utils.GetQueryRunner(tx, r.db)

	_, err := runner.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)
	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Where("is_deleted = FALSE").
		Scan(ctx)

	return user, err
}

func (r *Repository) RegiserUserDevice(ctx context.Context, userDevice *UserDevice) error {
	_, err := r.db.NewInsert().
		Model(userDevice).
		Exec(ctx)

	return err
}

func (r *Repository) FindActiveUserDeviceByUserId(ctx context.Context, userId string) (UserDevice, error) {
	userDevice := new(UserDevice)
	err := r.db.NewSelect().
		Model(userDevice).
		Where("user_id = ?", userId).
		Where("is_active = ?", true).
		Where("revoked_at = ?", nil).
		Limit(1).
		Scan(ctx)

	return *userDevice, err
}

func (r *Repository) UnRegisterUserDevice(ctx context.Context, deviceId string) error {
	_, err := r.db.NewUpdate().
		Model(&UserDevice{}).
		Set("revoked_at = ?", time.Now()).
		Set("updated_at = ?", time.Now()).
		Set("is_active = ?", false).
		Where("device_id = ?", deviceId).
		Exec(ctx)

	return err
}

func (r *Repository) GetActiveUserDeviceTokens(ctx context.Context, userId string) (string, error) {
	devices := new(UserDevice)

	err := r.db.NewSelect().
		Model(devices).
		Where("user_id = ?", userId).
		Where("is_active = TRUE").
		Where("is_deleted = FALSE").
		Where("revoked_at IS NULL").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return "", err
	}

	return devices.Token, nil
}

func (r *Repository) GetUserByEmailOrUsername(ctx context.Context, username, email string) ([]User, error) {
	users := make([]User, 0, 2)

	err := r.db.NewSelect().
		Model(&users).
		Where("username = ?", username).
		WhereOr("email = ?", email).
		Where("is_deleted = FALSE").
		Scan(ctx)

	return users, err
}
