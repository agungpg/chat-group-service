package files

import (
	"time"

	"github.com/uptrace/bun"
)

type FileScope string

const (
	FileScopeAvatar         FileScope = "avatar"
	FileScopeChatAttachment FileScope = "chat_attachment"
)

type Files struct {
	bun.BaseModel `bun:"table:files"`

	ID               string    `bun:"id,pk,notnul"`
	Bucket           string    `bun:"bucket,notnull"`
	ObjectKey        string    `bun:"object_key,notnull"`
	OriginalFilename string    `bun:"original_filename,notnull"`
	ContentType      string    `bun:"content_type,notnull"`
	SizeBytes        int64     `bun:"size_bytes,notnull"`
	Scope            FileScope `bun:"scope,notnull"`
	CreatedBy        string    `bun:"created_by,notnull"`
	CreatedAt        time.Time `bun:"created_at,notnull,default:current_timestamp"`
	DeletedAt        time.Time `bun:"deleted_at,nullzero"`
}
