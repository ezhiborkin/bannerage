package handler

import (
	"banners/domain/models"
	"banners/internal/errorwriter"
	"banners/internal/storage"
	"banners/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type BannerProvider interface {
	PostBanner(ctx context.Context, banner *models.Banner) (int, error)
	GetUserBanner(ctx context.Context, tagID int, featureID int) (*models.Banner, error)
	ChooseRevision(ctx context.Context, bannerID int, revisionID int) error
	ListRevisions(ctx context.Context, bannerID int, limit int, offset int) (*[]models.Banner, error)
	ListBanners(ctx context.Context, featureID int, tagID int, limit int, offset int) (*[]models.Banner, error)
	DeleteBanner(ctx context.Context, bannerID int) error
	DeleteUserBannerByFeatureTag(ctx context.Context, tagID int, featureID int) error
	PatchBanner(ctx context.Context, banner *models.Banner) error
	GetUserBannerCache(ctx context.Context, tagID int, featureID int) (*models.Banner, error)
	SetUserBannerCache(ctx context.Context, tagID int, featureID int, banner models.Banner) error
}

func (h *Handler) postBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.postBanner"

	log := h.log.With(slog.String("op", op))

	type postBanner struct {
		Content   json.RawMessage `json:"content"`
		FeatureID *int64          `json:"feature_id"`
		TagIDs    []int64         `json:"tag_ids"`
		IsActive  *bool           `json:"is_active"`
	}

	var bannerReq postBanner
	err := json.NewDecoder(r.Body).Decode(&bannerReq)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		errorwriter.WriteError(w, "failed to decode request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", sl.Err(err))
		errorwriter.WriteError(w, "empty request", http.StatusBadRequest)
		return
	}

	if bannerReq.Content == nil || bannerReq.FeatureID == nil || bannerReq.IsActive == nil || bannerReq.TagIDs == nil {
		log.Error("failed to create banner: missing required fields in request body")
		errorwriter.WriteError(w, "failed to create banner: missing required fields in request body", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	banner := &models.Banner{
		FeatureID: *bannerReq.FeatureID,
		TagIDs:    bannerReq.TagIDs,
		Content:   bannerReq.Content,
		IsActive:  *bannerReq.IsActive,
	}

	bannerID, err := h.bannerProvider.PostBanner(r.Context(), banner)
	if err != nil {
		log.Error("failed to create banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to create banner", http.StatusInternalServerError)
		return
	}

	type postUserBanner struct {
		Message  string `json:"message"`
		BannerID int    `json:"banner_id"`
	}

	response := postUserBanner{
		Message:  "Successfully created banner.",
		BannerID: bannerID,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		errorwriter.WriteError(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error("failed to create banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to create banner", http.StatusInternalServerError)
	}
}

func (h *Handler) getUserBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getUserBanner"

	log := h.log.With(slog.String("op", op))

	tagIDStr := r.URL.Query().Get("tag_id")
	featureIDStr := r.URL.Query().Get("feature_id")
	useLastRevStr := r.URL.Query().Get("use_last_revision")

	if tagIDStr == "" || featureIDStr == "" {
		log.Error("tagID or featureID is not provided")
		errorwriter.WriteError(w, "tagID or featureID is not provided", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		log.Error("tagID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "tagID is not a number", http.StatusBadRequest)
		return
	}

	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		log.Error("featureID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "featureID is not a number", http.StatusBadRequest)
		return
	}

	var useLastRev bool
	if useLastRevStr == "" {
		useLastRev = false
	} else {
		useLastRev, err = strconv.ParseBool(useLastRevStr)
		if err != nil {
			log.Error("useLastRev is not a bool", sl.Err(err))
			errorwriter.WriteError(w, "useLastRev is not a bool", http.StatusBadRequest)
			return
		}
	}

	var banner *models.Banner
	if useLastRev == false {
		banner, err = h.bannerProvider.GetUserBannerCache(r.Context(), tagID, featureID)
		if err != nil {
			log.Error("failed to get banner", sl.Err(err))
			errorwriter.WriteError(w, "failed to get banner", http.StatusNotFound)
			return
		}
	} else {
		banner, err = h.bannerProvider.GetUserBanner(r.Context(), tagID, featureID)
		if errors.Is(err, storage.ErrBannerNotFound) {
			log.Info("banner not found", sl.Err(err))
			errorwriter.WriteError(w, "banner not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Error("failed to get banner", sl.Err(err))
			errorwriter.WriteError(w, "failed to get banner", http.StatusNotFound)
			return
		}
	}

	if banner.IsActive == false {
		if r.Context().Value("role") != "admin" {
			log.Error("you are not admin")
			errorwriter.WriteError(w, "you are not admin", http.StatusUnauthorized)
			return
		}
	}

	type createBanner struct {
		Content json.RawMessage `json:"content"`
	}

	response := createBanner{
		Content: banner.Content,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		errorwriter.WriteError(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error("failed to get banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to get banner", http.StatusInternalServerError)
	}

}

func (h *Handler) chooseBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.chooseBanner"

	log := h.log.With(slog.String("op", op))

	bannerIDStr := r.URL.Query().Get("banner_id")
	revisionIDStr := r.URL.Query().Get("revision_id")

	if bannerIDStr == "" || revisionIDStr == "" {
		log.Error("bannerID or revisionID is not provided")
		errorwriter.WriteError(w, "bannerID or revisionID is not provided", http.StatusBadRequest)
		return
	}

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		log.Error("bannerID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "bannerID is not a number", http.StatusBadRequest)
		return
	}

	revisionID, err := strconv.Atoi(revisionIDStr)
	if err != nil {
		log.Error("revisionID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "revisionID is not a number", http.StatusBadRequest)
		return
	}

	err = h.bannerProvider.ChooseRevision(r.Context(), bannerID, revisionID)
	if errors.Is(err, storage.ErrFailedRevisionChange) {
		log.Info("failed to choose a revision", sl.Err(err))
		errorwriter.WriteError(w, "failed to choose a revision", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Error("failed to choose a version", sl.Err(err))
		errorwriter.WriteError(w, "failed to choose a version", http.StatusBadRequest)
		return
	}

	type chooseBanner struct {
		Message    string `json:"message"`
		BannerID   int    `json:"banner_id"`
		RevisionID int    `json:"revision_id"`
	}

	response := chooseBanner{
		Message:    "Successfully chosen a revision",
		RevisionID: revisionID,
		BannerID:   bannerID,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		errorwriter.WriteError(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error("failed to choose revision", sl.Err(err))
		errorwriter.WriteError(w, "failed to choose revision", http.StatusInternalServerError)
	}
}

func (h *Handler) listRevisions(w http.ResponseWriter, r *http.Request) {
	const op = "handler.listRevisions"

	log := h.log.With(slog.String("op", op))

	bannerIDStr := r.PathValue("banner_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr == "" {
		limitStr = "5"
	}

	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error("limit is not a number", sl.Err(err))
		errorwriter.WriteError(w, "limit is not a number", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Error("offset is not a number", sl.Err(err))
		errorwriter.WriteError(w, "offset is not a number", http.StatusBadRequest)
		return
	}

	if bannerIDStr == "" {
		log.Error("bannerID  is not provided")
		errorwriter.WriteError(w, "bannerID  is not provided", http.StatusBadRequest)
		return
	}

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		log.Error("bannerID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "bannerID is not a number", http.StatusBadRequest)
		return
	}

	if limit < 0 || limit > 100 {
		log.Error("limit is out of range")
		errorwriter.WriteError(w, "limit is out of range", http.StatusBadRequest)
		return
	}

	if offset < 0 {
		log.Error("offset is out of range")
		errorwriter.WriteError(w, "offset is out of range", http.StatusBadRequest)
		return
	}

	revisions, err := h.bannerProvider.ListRevisions(r.Context(), bannerID, limit, offset)
	if err != nil {
		log.Error("failed to list revisions", sl.Err(err))
		errorwriter.WriteError(w, "failed to list revisions", http.StatusInternalServerError)
		return
	}

	if len(*revisions) == 0 {
		log.Error("no revisions found")
		errorwriter.WriteError(w, "no revisions found", http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(*revisions)
	if err != nil {
		log.Error("failed to list revisions", sl.Err(err))
		errorwriter.WriteError(w, "failed to list revisions", http.StatusInternalServerError)
	}
}

func (h *Handler) listBanners(w http.ResponseWriter, r *http.Request) {
	const op = "handler.listBanners"

	log := h.log.With(slog.String("op", op))

	tagIDStr := r.URL.Query().Get("tag_id")
	featureIDStr := r.URL.Query().Get("feature_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if tagIDStr == "" || featureIDStr == "" {
		log.Error("tagID or featureID is not provided")
		errorwriter.WriteError(w, "tagID or featureID is not provided", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		log.Error("tagID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "tagID is not a number", http.StatusBadRequest)
		return
	}

	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		log.Error("featureID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "featureID is not a number", http.StatusBadRequest)
		return
	}

	if limitStr == "" {
		limitStr = "5"
	}

	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error("limit is not a number", sl.Err(err))
		errorwriter.WriteError(w, "limit is not a number", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Error("offset is not a number", sl.Err(err))
		errorwriter.WriteError(w, "offset is not a number", http.StatusBadRequest)
		return
	}

	if limit < 0 || limit > 100 {
		log.Error("limit is out of range")
		errorwriter.WriteError(w, "limit is out of range", http.StatusBadRequest)
		return
	}

	if offset < 0 {
		log.Error("offset is out of range")
		errorwriter.WriteError(w, "offset is out of range", http.StatusBadRequest)
		return
	}

	banners, err := h.bannerProvider.ListBanners(r.Context(), featureID, tagID, limit, offset)
	if err != nil {
		log.Error("failed to list banners", sl.Err(err))
		errorwriter.WriteError(w, "failed to list banners", http.StatusInternalServerError)
		return
	}

	fmt.Println(len(*banners))
	if len(*banners) == 0 {
		errorwriter.WriteError(w, "no banners found", http.StatusNoContent)
		w.Write([]byte("no banners found"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(*banners)
	if err != nil {
		log.Error("failed to list banners", sl.Err(err))
		errorwriter.WriteError(w, "failed to list banners", http.StatusInternalServerError)
	}
}

func (h *Handler) patchBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.patchBanner"

	log := h.log.With(slog.String("op", op))

	bannerIDStr := r.PathValue("id")

	if bannerIDStr == "" {
		log.Error("bannerID  is not provided")
		errorwriter.WriteError(w, "bannerID  is not provided", http.StatusBadRequest)
		return
	}

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		log.Error("bannerID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "bannerID is not a number", http.StatusBadRequest)
		return
	}

	banner := &models.Banner{}
	err = json.NewDecoder(r.Body).Decode(banner)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		errorwriter.WriteError(w, "failed to decode request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", sl.Err(err))
		errorwriter.WriteError(w, "empty request", http.StatusBadRequest)
		return
	}

	banner.BannerID = int64(bannerID)

	//TODO

	log.Info("request body decoded")

	err = h.bannerProvider.PatchBanner(r.Context(), banner)
	if err != nil {
		log.Error("failed to patch banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to patch banner", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("patched banner: %v", bannerID)))
	if err != nil {
		log.Error("failed to patch a banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to patch a banner", http.StatusInternalServerError)
	}
}

func (h *Handler) deleteBanner(w http.ResponseWriter, r *http.Request) {
	const op = "storage.postgresql.deleteBanner"

	log := h.log.With(slog.String("op", op))

	bannerIDStr := r.PathValue("id")

	if bannerIDStr == "" {
		log.Error("bannerID  is not provided")
		errorwriter.WriteError(w, "bannerID  is not provided", http.StatusBadRequest)
		return
	}

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		log.Error("bannerID is not a number", sl.Err(err))
		errorwriter.WriteError(w, "bannerID is not a number", http.StatusBadRequest)
		return
	}

	err = h.bannerProvider.DeleteBanner(r.Context(), bannerID)
	if err != nil {
		log.Error("failed to delete banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to delete banner", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("deleted banner: %v", bannerID)))
	if err != nil {
		log.Error("failed to delete banner", sl.Err(err))
		errorwriter.WriteError(w, "failed to delete banner", http.StatusInternalServerError)
	}
}

func (h *Handler) deleteBannerFeatureTag(deleteCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "storage.postgresql.deleteBannerFeatureTag"

		log := h.log.With(slog.String("op", op))

		tagIDStr := r.URL.Query().Get("tag_id")
		featureIDStr := r.URL.Query().Get("feature_id")

		if tagIDStr == "" || featureIDStr == "" {
			log.Error("tagID or featureID is not provided")
			errorwriter.WriteError(w, "tagID or featureID is not provided", http.StatusBadRequest)
			return
		}

		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			log.Error("tagID is not a number", sl.Err(err))
			errorwriter.WriteError(w, "tagID is not a number", http.StatusBadRequest)
			return
		}

		featureID, err := strconv.Atoi(featureIDStr)
		if err != nil {
			log.Error("featureID is not a number", sl.Err(err))
			errorwriter.WriteError(w, "featureID is not a number", http.StatusBadRequest)
			return
		}

		err = h.bannerProvider.DeleteUserBannerByFeatureTag(deleteCtx, tagID, featureID)
		if err != nil {
			log.Error("failed to delete banner", sl.Err(err))
			errorwriter.WriteError(w, "failed to delete banner", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("deleted banner with featureID: %d, tagID: %d, OR NOT)))", featureID, tagID)))
		if err != nil {
			log.Error("failed to delete banner", sl.Err(err))
			errorwriter.WriteError(w, "failed to delete banner", http.StatusInternalServerError)
		}

	}
}
