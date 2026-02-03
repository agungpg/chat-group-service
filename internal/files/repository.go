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
		Model(&file).
		Exec(ctx)

	return err
}
