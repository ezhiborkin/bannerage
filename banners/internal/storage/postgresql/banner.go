package postgresql

import (
	"banners/domain/models"
	"banners/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strconv"
	"strings"
)

func (s *Storage) GetUsersBannerStorage(ctx context.Context, tagID int, featureID int) (*models.Banner, error) {
	const op = "storage.postgresql.GetUsersBanner"

	query, args, err := sq.Select("br.content, br.is_active").
		From("banner_revisions br").
		Join("revision_tags rt ON br.revision_id = rt.revision_id").
		Where(sq.And{
			sq.Eq{"rt.tag_id": tagID},
			sq.Eq{"br.feature_id": featureID},
		}).
		Where(sq.Expr("br.banner_id IN (SELECT banner_id FROM banners WHERE chosen_revision_id = br.revision_id)")).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	var banner models.Banner
	err = row.Scan(&banner.Content, &banner.IsActive)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrBannerNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &banner, nil
}

func (s *Storage) PostBannerStorage(ctx context.Context, banner *models.Banner) (int, error) {
	const op = "storage.postgresql.PostBanner"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	bannerInsert := sq.Insert("banners").
		Values(sq.Expr("DEFAULT")).
		Suffix("RETURNING banner_id")

	var bannerID int
	err = bannerInsert.RunWith(tx).PlaceholderFormat(sq.Dollar).ScanContext(ctx, &bannerID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	bannerRevInsert := sq.Insert("banner_revisions").
		Columns("banner_id", "feature_id", "content").
		Values(bannerID, banner.FeatureID, banner.Content).
		Suffix("RETURNING revision_id")

	var revisionID int
	err = bannerRevInsert.RunWith(tx).PlaceholderFormat(sq.Dollar).ScanContext(ctx, &revisionID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	bannerTagsInsert := sq.Insert("revision_tags").
		Columns("revision_id", "tag_id")

	for _, tagID := range banner.TagIDs {
		bannerTagsInsert = bannerTagsInsert.Values(revisionID, tagID)
	}

	_, err = bannerTagsInsert.RunWith(tx).PlaceholderFormat(sq.Dollar).ExecContext(ctx)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	bannerUpdate, args, err := sq.Update("banners").
		Set("chosen_revision_id", revisionID).
		Where(sq.Eq{"banner_id": bannerID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, bannerUpdate, args...)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return bannerID, nil
}

func (s *Storage) ChooseRevisionStorage(ctx context.Context, bannerID int, revisionID int) error {
	const op = "storage.postgresql.ChooseRevision"

	checkQuery := sq.Select("*").
		From("banner_revisions").
		Where(sq.Eq{"banner_id": bannerID, "revision_id": revisionID}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	err := checkQuery.RunWith(s.db).QueryRowContext(ctx).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", op, err)
	}

	chosenRevInsert := sq.Update("banners").
		Set("chosen_revision_id", revisionID).
		Where(sq.Eq{"banner_id": bannerID}).
		Suffix("RETURNING chosen_revision_id")

	err = chosenRevInsert.RunWith(s.db).PlaceholderFormat(sq.Dollar).QueryRowContext(ctx).Scan(&revisionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ListRevisionsStorage(ctx context.Context, bannerID int, limit int, offset int) (*[]models.Banner, error) {
	const op = "storage.postgresql.ListRevisionsStorage"

	query, args, err := sq.Select("br.*, ARRAY_TO_STRING(ARRAY_AGG(rt.tag_id), ', ') AS tags").
		From("banner_revisions br").
		Join("revision_tags rt ON br.revision_id = rt.revision_id").
		Where(sq.Eq{"banner_id": bannerID}).
		GroupBy("br.revision_id, br.banner_id, br.feature_id, br.is_active, br.content, br.created_at, br.updated_at").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var revisions []models.Banner
	var tagIDsStr string
	for rows.Next() {
		var revision models.Banner
		if err = rows.Scan(&revision.Revision, &revision.BannerID, &revision.FeatureID, &revision.IsActive, &revision.Content, &revision.CreatedAt, &revision.UpdatedAT, &tagIDsStr); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		revision.TagIDs, err = parseTagIDs(tagIDsStr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		revisions = append(revisions, revision)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &revisions, nil
}

func (s *Storage) ListBannersStorage(ctx context.Context, featureID int, tagID int, limit int, offset int) (*[]models.Banner, error) {
	const op = "storage.postgresql.ListBannersStorage"
	//subQuery := sq.Select("br.revision_id").
	//	From("banners b").
	//	Join("banner_revisions br ON b.chosen_revision_id = br.revision_id").
	//	Where(sq.Eq{"br.feature_id": featureID}).
	//	Where("br.revision_id = b.chosen_revision_id").
	//	GroupBy("br.revision_id").
	//	PlaceholderFormat(sq.Dollar)
	//
	//mainQuery, args, err := sq.
	//	Select("b.banner_id", "br.feature_id", "br.is_active", "br.created_at", "br.updated_at", "br.revision_id", "br.content", "ARRAY_TO_STRING(ARRAY_AGG(rt.tag_id), ', ') AS tag_ids").
	//	From("banners b").
	//	Join("banner_revisions br ON b.chosen_revision_id = br.revision_id").
	//	Join("revision_tags rt ON br.revision_id = rt.revision_id").
	//	Where(sq.Expr("br.revision_id IN (?)", subQuery)).
	//	Where(sq.Eq{"rt.tag_id": tagID}).
	//	GroupBy("b.banner_id", "br.feature_id", "br.is_active", "br.created_at", "br.updated_at", "br.revision_id", "br.content").
	//	Offset(uint64(offset)).
	//	Limit(uint64(limit)).
	//	PlaceholderFormat(sq.Dollar).
	//	ToSql()
	query, args, err := sq.Select("b.banner_id", "br.feature_id", "br.is_active", "br.created_at", "br.updated_at", "br.revision_id", "br.content", "ARRAY_TO_STRING(ARRAY_AGG(rt.tag_id), ', ') AS tag_ids").
		From("banners b").
		Join("banner_revisions br ON b.chosen_revision_id = br.revision_id").
		Join("revision_tags rt ON br.revision_id = rt.revision_id").
		Where(sq.Eq{"br.feature_id": featureID}).
		Where(sq.Expr("rt.revision_id IN (SELECT revision_id FROM revision_tags WHERE tag_id = ?)", tagID)).
		Where(sq.Expr("br.revision_id IN (?)",
			sq.Select("br.revision_id").
				From("banners b").
				Join("banner_revisions br ON b.chosen_revision_id = br.revision_id").
				Where(sq.Eq{"br.feature_id": featureID}).
				Where(sq.Expr("br.revision_id = b.chosen_revision_id")).
				GroupBy("br.revision_id"),
		)).
		GroupBy("b.banner_id", "br.feature_id", "br.is_active", "br.created_at", "br.updated_at", "br.revision_id", "br.content").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var banners []models.Banner
	var tagIDsStr string
	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner.BannerID, &banner.FeatureID, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAT, &banner.Revision, &banner.Content, &tagIDsStr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		banner.TagIDs, err = parseTagIDs(tagIDsStr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		banners = append(banners, banner)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &banners, nil
}

func (s *Storage) PatchBannerStorage(ctx context.Context, banner *models.Banner) error {
	const op = "storage.postgresql.PatchBannerStorage"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	subQuery := sq.Select("chosen_revision_id").
		From("banners").
		Where(sq.Eq{"banner_id": banner.BannerID}).
		Limit(1)

	var chosenRevisionID int
	err = subQuery.RunWith(tx).PlaceholderFormat(sq.Dollar).QueryRowContext(ctx).Scan(&chosenRevisionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	updateBuilder := sq.Update("banner_revisions").
		Where(sq.Eq{"revision_id": chosenRevisionID}).
		Set("is_active", banner.IsActive)

	if banner.FeatureID != 0 {
		updateBuilder = updateBuilder.Set("feature_id", banner.FeatureID)
	}

	if banner.Content != nil {
		updateBuilder = updateBuilder.Set("content", banner.Content)
	}

	if len(banner.TagIDs) != 0 {
		query, args, err := sq.Delete("revision_tags").
			Where(sq.Eq{"revision_id": chosenRevisionID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		bannerTagsInsert := sq.Insert("revision_tags").
			Columns("revision_id", "tag_id")

		for _, tagID := range banner.TagIDs {
			bannerTagsInsert = bannerTagsInsert.Values(chosenRevisionID, tagID)
		}

		_, err = bannerTagsInsert.RunWith(tx).PlaceholderFormat(sq.Dollar).ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	query, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteBannerStorage(ctx context.Context, bannerID int) error {
	const op = "storage.postgresql.DeleteBannerStorage"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query, args, err := sq.Delete("banners").
		Where(sq.Eq{"banner_id": bannerID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserBannerByFeatureTagStorage(ctx context.Context, tagID int, featureID int) error {
	const op = "storage.postgresql.DeleteUserBannerByFeatureTagStorage"
	task := func() error {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()

		query, args, err := sq.Delete("banners").
			Where(sq.Expr("chosen_revision_id IN (SELECT br.revision_id FROM banner_revisions br JOIN revision_tags rt ON br.revision_id = rt.revision_id WHERE br.feature_id = ? AND rt.tag_id = ?)", featureID, tagID)).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}

	s.workerPool.AddTask(task)

	return nil
}

func parseTagIDs(tagIDsStr string) ([]int64, error) {
	var tagIDs []int64
	tags := strings.Split(tagIDsStr, ",")
	for _, tag := range tags {
		tag = strings.TrimSpace(tag) // Удаление пробелов
		tagID, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}
	return tagIDs, nil
}

func parseTagIDsBanner(tagIDsStr string) ([]int64, error) {
	if len(tagIDsStr) < 3 { // Проверяем, что строка содержит хотя бы "{x}"
		return nil, errors.New("invalid tagIDs string format")
	}
	tagIDsStr = tagIDsStr[1 : len(tagIDsStr)-1] // Удаляем круглые скобки в начале и конце
	tagIDStrs := strings.Split(tagIDsStr, ",")
	tagIDs := make([]int64, len(tagIDStrs))
	for i, tagIDStr := range tagIDStrs {
		tagID, err := strconv.ParseInt(strings.TrimSpace(tagIDStr), 10, 64)
		if err != nil {
			return nil, err
		}
		tagIDs[i] = tagID
	}
	return tagIDs, nil
}
