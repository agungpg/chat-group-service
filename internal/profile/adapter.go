package profile

import (
	"context"

	"github.com/agungpg/group-chat-service/internal/files"
)

type FileAdapter interface {
	GetPresignView(ctx context.Context, fileId string) (*files.UploadPresignResponseDTO, error)
}
