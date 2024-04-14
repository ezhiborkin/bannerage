package models

import (
	"encoding/json"
	"time"
)

type Banner struct {
	BannerID  int64           `json:"banner_id,omitempty"`
	TagIDs    []int64         `json:"tag_ids,omitempty"`
	FeatureID int64           `json:"feature_id,omitempty"`
	Revision  int64           `json:"revision_id,omitempty"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAT time.Time       `json:"updated_at,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
}
