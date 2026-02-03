package files

import (
	"context"

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

func (r *Repository) CreateFile(ctx context.Context, file *Files) error {

	_, err := r.db.NewInsert().
		Model(file).
		Exec(ctx)

	return err
}

func (r *Repository) GetFileById(ctx context.Context, fileId string) (*Files, error) {
	file := &Files{}
	err := r.db.NewSelect().
		Where("deleted_at IS NULL").
		Where("status <> ?", FileStatusDeleted).
		Where("id = ?", fileId).
		Model(file).
		Scan(ctx)

	return file, err
}

func (r *Repository) UpdateFileStatus(ctx context.Context, fileId string, status FileStatus) error {

	_, err := r.db.NewUpdate().
		Model(&Files{}).
		Set("status = ?", status).
		Where("id = ?", fileId).
		Exec(ctx)

	return err
}
