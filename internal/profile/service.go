package profile

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetProfileByUserId(ctx context.Context, userId string) (*ProfileResponse, error) {
	uf, err := s.repo.GetProfileByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	profile := &ProfileResponse{
		UserID:      uf.UserID,
		DisplayName: uf.DisplayName,
		AvatarURL:   uf.AvatarURL,
		Bio:         uf.Bio,
	}
	return profile, err
}
