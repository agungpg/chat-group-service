package files

import "time"

type UploadPresignRequestDTO struct {
	FileName    string `json:"filename"`
	SizeBytes   int64  `json:"sizeInBytes"`
	ContentType string `json:"contentType"`
	FileType    string `json:"fileType"`
}

type UploadPresignResponseDTO struct {
	PresignedUrl string        `json:"presignedUrl"`
	ExpiresIn    time.Duration `json:"expiresIn"`
	FileId       string        `json:"fileId"`
}
