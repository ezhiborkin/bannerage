package storage

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrBannerNotFound       = errors.New("banner not found")
	ErrFailedRevisionChange = errors.New("failed to choose a revision")
	ErrRevisionDoesNotExist = errors.New("chosen revision does not exist for this banner")
	ErrNotFoundInCache      = errors.New("value not found in cache")
)
