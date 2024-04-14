package service

import (
	"banners/domain/models"
	"banners/internal/storage"
	"banners/lib/logger/sl"
	"context"
	"errors"
	"fmt"
)

type BannerStorage interface {
	PostBannerStorage(ctx context.Context, banner *models.Banner) (int, error)
	GetUsersBannerStorage(ctx context.Context, tagID int, featureID int) (*models.Banner, error)
	ChooseRevisionStorage(ctx context.Context, bannerID int, revisionID int) error
	ListRevisionsStorage(ctx context.Context, bannerID int, limit int, offset int) (*[]models.Banner, error)
	ListBannersStorage(ctx context.Context, featureID int, tagID int, limit int, offset int) (*[]models.Banner, error)
	DeleteBannerStorage(ctx context.Context, bannerID int) error
	DeleteUserBannerByFeatureTagStorage(ctx context.Context, tagID int, featureID int) error
	PatchBannerStorage(ctx context.Context, banner *models.Banner) error
}

func (s *Service) PostBanner(ctx context.Context, banner *models.Banner) (int, error) {
	const op = "service.PostBanner"

	bannerID, err := s.bannerStorage.PostBannerStorage(ctx, banner)
	if err != nil {
		s.log.Error("failed to post banner", sl.Err(err))

		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return bannerID, nil
}

func (s *Service) GetUserBanner(ctx context.Context, tagID int, featureID int) (*models.Banner, error) {
	const op = "service.GetUserBanner"

	banner, err := s.bannerStorage.GetUsersBannerStorage(ctx, tagID, featureID)
	if err != nil {
		s.log.Error("failed to get user banner", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return banner, nil
}

func (s *Service) GetUserBannerCache(ctx context.Context, tagID int, featureID int) (*models.Banner, error) {
	const op = "service.GetUserBannerCache"

	key := fmt.Sprintf("banner:%d:%d", tagID, featureID)
	banner, err := s.c.Get(ctx, key)
	if errors.Is(err, storage.ErrNotFoundInCache) {
		banner, err := s.bannerStorage.GetUsersBannerStorage(ctx, tagID, featureID)
		if err != nil {
			s.log.Error("failed to get user banner", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = s.SetUserBannerCache(ctx, tagID, featureID, *banner)
		if err != nil {
			s.log.Error("failed to set user banner in cache", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return banner, nil
	}
	if err != nil {
		s.log.Error("failed to get user banner from cache", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//var content json.RawMessage
	//err = content.UnmarshalJSON([]byte(bannerContentStr))
	//if err != nil {
	//	s.log.Error("failed to unmarshal content string to json", sl.Err(err))
	//	return nil, fmt.Errorf("%s: %w", op, err)
	//}

	return banner, nil
}

func (s *Service) SetUserBannerCache(ctx context.Context, tagID int, featureID int, banner models.Banner) error {
	const op = "service.SetUserBannerCache"

	//bytes, err := content.MarshalJSON()
	//if err != nil {
	//	s.log.Error("failed to marshal content to string", sl.Err(err))
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	key := fmt.Sprintf("banner:%d:%d", tagID, featureID)
	err := s.c.Set(ctx, key, banner)
	if err != nil {
		s.log.Error("failed to set user banner in cache", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ChooseRevision(ctx context.Context, bannerID int, revisionID int) error {
	const op = "service.ChooseRevision"

	err := s.bannerStorage.ChooseRevisionStorage(ctx, bannerID, revisionID)
	if err != nil {
		s.log.Error("failed to choose revision", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ListRevisions(ctx context.Context, bannerID int, limit int, offset int) (*[]models.Banner, error) {
	const op = "service.ListRevisions"

	revisions, err := s.bannerStorage.ListRevisionsStorage(ctx, bannerID, limit, offset)
	if err != nil {
		s.log.Error("failed to list revisions", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return revisions, nil
}

func (s *Service) ListBanners(ctx context.Context, featureID int, tagID int, limit int, offset int) (*[]models.Banner, error) {
	const op = "service.ListBanners"

	banners, err := s.bannerStorage.ListBannersStorage(ctx, featureID, tagID, limit, offset)
	if err != nil {
		s.log.Error("failed to list banners", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return banners, nil
}

func (s *Service) DeleteBanner(ctx context.Context, bannerID int) error {
	const op = "service.DeleteBanner"

	err := s.bannerStorage.DeleteBannerStorage(ctx, bannerID)
	if err != nil {
		s.log.Error("failed to list revisions", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteUserBannerByFeatureTag(ctx context.Context, tagID int, featureID int) error {
	const op = "service.DeleteUserBannerByFeatureTag"

	err := s.bannerStorage.DeleteUserBannerByFeatureTagStorage(ctx, tagID, featureID)
	if err != nil {
		s.log.Error("failed to list revisions", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) PatchBanner(ctx context.Context, banner *models.Banner) error {
	const op = "service.PatchBanner"

	err := s.bannerStorage.PatchBannerStorage(ctx, banner)
	if err != nil {
		s.log.Error("failed to list revisions", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
