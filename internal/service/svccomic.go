package service

import (
	"context"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

//
// Comic
//

func (svc Service) AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.Titles = []*model.ComicTitle{}
		v.Covers = []*model.ComicCover{}
		v.Synopses = []*model.ComicSynopsis{}
		v.Chapters = []*model.ComicChapter{}
		v.Externals = []*model.ComicExternal{}
		v.Categories = []*model.Category{}
		v.Tags = []*model.Tag{}
		v.Relations = []*model.ComicRelation{}
	}

	return svc.database.AddComic(ctx, data, v)
}

func (svc Service) GetComicByCode(ctx context.Context, code string) (*model.Comic, error) {
	result, err := svc.database.GetComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	})
	if err != nil {
		return nil, err
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		titles, err := svc.database.ListComicTitle(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Titles = titles
		return nil
	})
	g.Go(func() error {
		covers, err := svc.database.ListComicCover(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Covers = covers
		return nil
	})
	g.Go(func() error {
		synopses, err := svc.database.ListComicSynopsis(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Synopses = synopses
		return nil
	})
	g.Go(func() error {
		chapters, err := svc.database.ListComicChapter(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Chapters = chapters
		return nil
	})
	g.Go(func() error {
		externals, err := svc.database.ListComicExternal(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Externals = externals
		return nil
	})
	g.Go(func() error {
		categories0, err := svc.database.ListComicCategory(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		if len(categories0) == 0 {
			result.Categories = []*model.Category{}
			return nil
		}

		conditions := make([]any, len(categories0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for _, category := range categories0 {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: category.CategoryID,
			})
		}

		categories1, err := svc.database.ListCategory(gctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Categories = categories1
		return nil
	})
	g.Go(func() error {
		tags0, err := svc.database.ListComicTag(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		if len(tags0) == 0 {
			result.Tags = []*model.Tag{}
			return nil
		}

		conditions := make([]any, len(tags0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for _, tag := range tags0 {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: tag.TagID,
			})
		}

		tags1, err := svc.database.ListTag(gctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Tags = tags1
		return nil
	})
	g.Go(func() error {
		relations, err := svc.database.ListComicRelation(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicRelationParentID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Relations = relations
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) UpdateComicByCode(ctx context.Context, code string, data model.SetComic, v *model.Comic) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if err := svc.database.UpdateComic(ctx, data, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	}, v); err != nil {
		return err
	}

	if v != nil {
		g, gctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			titles, err := svc.database.ListComicTitle(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Titles = titles
			return nil
		})
		g.Go(func() error {
			covers, err := svc.database.ListComicCover(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Covers = covers
			return nil
		})
		g.Go(func() error {
			synopses, err := svc.database.ListComicSynopsis(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Synopses = synopses
			return nil
		})
		g.Go(func() error {
			chapters, err := svc.database.ListComicChapter(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Chapters = chapters
			return nil
		})
		g.Go(func() error {
			externals, err := svc.database.ListComicExternal(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Externals = externals
			return nil
		})
		g.Go(func() error {
			categories0, err := svc.database.ListComicCategory(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			if len(categories0) == 0 {
				v.Categories = []*model.Category{}
				return nil
			}

			conditions := make([]any, len(categories0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for _, category := range categories0 {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: category.CategoryID,
				})
			}

			categories1, err := svc.database.ListCategory(gctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Categories = categories1
			return nil
		})
		g.Go(func() error {
			tags0, err := svc.database.ListComicTag(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			if len(tags0) == 0 {
				v.Tags = []*model.Tag{}
				return nil
			}

			conditions := make([]any, len(tags0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for _, tag := range tags0 {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: tag.TagID,
				})
			}

			tags1, err := svc.database.ListTag(gctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Tags = tags1
			return nil
		})
		g.Go(func() error {
			relations, err := svc.database.ListComicRelation(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicRelationParentID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Relations = relations
			return nil
		})
		if err := g.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (svc Service) DeleteComicByCode(ctx context.Context, code string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic")
	}

	return svc.database.DeleteComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	}, nil)
}

func (svc Service) ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicPaginationMax {
			pagination.Limit = model.ComicPaginationMax
		}
	}

	result, err := svc.database.ListComic(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBComicGenericComicID,
				Value: r.ID,
			})
		}
		g, gctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			titles, err := svc.database.ListComicTitle(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Titles = []*model.ComicTitle{}
			}
			for _, title := range titles {
				for _, r := range result {
					if r.ID == title.ComicID {
						r.Titles = append(r.Titles, title)
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			covers, err := svc.database.ListComicCover(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Covers = []*model.ComicCover{}
			}
			for _, cover := range covers {
				for _, r := range result {
					if r.ID == cover.ComicID {
						r.Covers = append(r.Covers, cover)
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			synopses, err := svc.database.ListComicSynopsis(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Synopses = []*model.ComicSynopsis{}
			}
			for _, synopsis := range synopses {
				for _, r := range result {
					if r.ID == synopsis.ComicID {
						r.Synopses = append(r.Synopses, synopsis)
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			chapters, err := svc.database.ListComicChapter(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Chapters = []*model.ComicChapter{}
			}
			for _, chapter := range chapters {
				for _, r := range result {
					if r.ID == chapter.ComicID {
						r.Chapters = append(r.Chapters, chapter)
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			externals, err := svc.database.ListComicExternal(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Externals = []*model.ComicExternal{}
			}
			for _, external := range externals {
				for _, r := range result {
					if r.ID == external.ComicID {
						r.Externals = append(r.Externals, external)
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			categories0, err := svc.database.ListComicCategory(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			categories := map[uint]*model.Category{}
			for _, category := range categories0 {
				categories[category.CategoryID] = nil
			}
			conditions := make([]any, len(categories0)+1)
			for id := range categories {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: id,
				})
			}
			categories1, err := svc.database.ListCategory(gctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, category := range categories1 {
				categories[category.ID] = category
			}
			for _, r := range result {
				r.Categories = []*model.Category{}
			}
			for _, category := range categories0 {
				for _, r := range result {
					if r.ID == category.ComicID {
						r.Categories = append(r.Categories, categories[category.CategoryID])
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			tags0, err := svc.database.ListComicTag(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			tags := map[uint]*model.Tag{}
			for _, tag := range tags0 {
				tags[tag.TagID] = nil
			}
			conditions := make([]any, len(tags0)+1)
			for id := range tags {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: id,
				})
			}
			tags1, err := svc.database.ListTag(gctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, tag := range tags1 {
				tags[tag.ID] = tag
			}
			for _, r := range result {
				r.Tags = []*model.Tag{}
			}
			for _, tag := range tags0 {
				for _, r := range result {
					if r.ID == tag.ComicID {
						r.Tags = append(r.Tags, tags[tag.TagID])
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			conds := make([]any, len(result)+1)
			conds = append(conds, model.DBLogicalOR{})
			for _, r := range result {
				conds = append(conds, model.DBConditionalKV{
					Key:   model.DBComicRelationParentID,
					Value: r.ID,
				})
			}
			relations, err := svc.database.ListComicRelation(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Relations = []*model.ComicRelation{}
			}
			for _, relation := range relations {
				for _, r := range result {
					if r.ID == relation.ParentID {
						r.Relations = append(r.Relations, relation)
					}
				}
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (svc Service) CountComic(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComic(ctx, conds)
}

func (svc Service) ExistsComicByCode(ctx context.Context, code string) (bool, error) {
	return svc.database.ExistsComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	})
}

// Comic Title

func (svc Service) AddComicTitle(ctx context.Context, data model.AddComicTitle, v *model.ComicTitle) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic title")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicTitle(ctx, data, v)
}

func (svc Service) GetComicTitleBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicTitle, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.GetComicTitle(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	})
}

func (svc Service) UpdateComicTitleBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicTitle, v *model.ComicTitle) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic title")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.UpdateComicTitle(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, v)
}

func (svc Service) DeleteComicTitleBySID(ctx context.Context, sid model.ComicGenericSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic title")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.DeleteComicTitle(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, nil)
}

func (svc Service) ListComicTitle(ctx context.Context, params model.ListParams) ([]*model.ComicTitle, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicTitleOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicTitleOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicTitleOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicTitlePaginationMax {
			pagination.Limit = model.ComicTitlePaginationMax
		}
	}

	return svc.database.ListComicTitle(ctx, params)
}

func (svc Service) CountComicTitle(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicTitle(ctx, conds)
}

// Comic Cover

func (svc Service) AddComicCover(ctx context.Context, data model.AddComicCover, v *model.ComicCover) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic cover")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicCover(ctx, data, v)
}

func (svc Service) GetComicCoverBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicCover, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.GetComicCover(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	})
}

func (svc Service) UpdateComicCoverBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicCover, v *model.ComicCover) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic cover")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.UpdateComicCover(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, v)
}

func (svc Service) DeleteComicCoverBySID(ctx context.Context, sid model.ComicGenericSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic cover")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.DeleteComicCover(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, nil)
}

func (svc Service) ListComicCover(ctx context.Context, params model.ListParams) ([]*model.ComicCover, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicCoverOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicCoverOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicCoverOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicCoverPaginationMax {
			pagination.Limit = model.ComicCoverPaginationMax
		}
	}

	return svc.database.ListComicCover(ctx, params)
}

func (svc Service) CountComicCover(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicCover(ctx, conds)
}

// Comic Synopsis

func (svc Service) AddComicSynopsis(ctx context.Context, data model.AddComicSynopsis, v *model.ComicSynopsis) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic synopsis")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicSynopsis(ctx, data, v)
}

func (svc Service) GetComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicSynopsis, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.GetComicSynopsis(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	})
}

func (svc Service) UpdateComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicSynopsis, v *model.ComicSynopsis) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic synopsis")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.UpdateComicSynopsis(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, v)
}

func (svc Service) DeleteComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic synopsis")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.DeleteComicSynopsis(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, nil)
}

func (svc Service) ListComicSynopsis(ctx context.Context, params model.ListParams) ([]*model.ComicSynopsis, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicSynopsisOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicSynopsisOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicSynopsisOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicSynopsisPaginationMax {
			pagination.Limit = model.ComicSynopsisPaginationMax
		}
	}

	return svc.database.ListComicSynopsis(ctx, params)
}

func (svc Service) CountComicSynopsis(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicSynopsis(ctx, conds)
}

// Comic External

func (svc Service) AddComicExternal(ctx context.Context, data model.AddComicExternal, v *model.ComicExternal) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic external")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicExternal(ctx, data, v)
}

func (svc Service) GetComicExternalBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicExternal, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.GetComicExternal(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	})
}

func (svc Service) UpdateComicExternalBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicExternal, v *model.ComicExternal) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic external")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.UpdateComicExternal(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, v)
}

func (svc Service) DeleteComicExternalBySID(ctx context.Context, sid model.ComicGenericSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic external")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	return svc.database.DeleteComicExternal(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicGenericRID:     sid.RID,
	}, nil)
}

func (svc Service) ListComicExternal(ctx context.Context, params model.ListParams) ([]*model.ComicExternal, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicExternalOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicExternalOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicExternalOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicExternalPaginationMax {
			pagination.Limit = model.ComicExternalPaginationMax
		}
	}

	return svc.database.ListComicExternal(ctx, params)
}

func (svc Service) CountComicExternal(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicExternal(ctx, conds)
}

// Comic Category

func (svc Service) AddComicCategory(ctx context.Context, data model.AddComicCategory, v *model.ComicCategory) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic category")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicCategory(ctx, data, v)
}

func (svc Service) GetComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID) (*model.ComicCategory, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var categoryID any
	switch {
	case sid.CategoryID != nil:
		categoryID = sid.CategoryID
	case sid.CategorySID != nil:
		categoryID = model.DBCategorySIDToID(*sid.CategorySID)
	}
	return svc.database.GetComicCategory(ctx, map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBCategoryGenericCategoryID: categoryID,
	})
}

func (svc Service) UpdateComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID, data model.SetComicCategory, v *model.ComicCategory) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic category")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var categoryID any
	switch {
	case sid.CategoryID != nil:
		categoryID = sid.CategoryID
	case sid.CategorySID != nil:
		categoryID = model.DBCategorySIDToID(*sid.CategorySID)
	}
	return svc.database.UpdateComicCategory(ctx, data, map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBCategoryGenericCategoryID: categoryID,
	}, v)
}

func (svc Service) DeleteComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic category")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var categoryID any
	switch {
	case sid.CategoryID != nil:
		categoryID = sid.CategoryID
	case sid.CategorySID != nil:
		categoryID = model.DBCategorySIDToID(*sid.CategorySID)
	}
	return svc.database.DeleteComicCategory(ctx, map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBCategoryGenericCategoryID: categoryID,
	}, nil)
}

func (svc Service) ListComicCategory(ctx context.Context, params model.ListParams) ([]*model.ComicCategory, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicCategoryOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicCategoryOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicCategoryOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicCategoryPaginationMax {
			pagination.Limit = model.ComicCategoryPaginationMax
		}
	}

	return svc.database.ListComicCategory(ctx, params)
}

func (svc Service) CountComicCategory(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicCategory(ctx, conds)
}

// Comic Tag

func (svc Service) AddComicTag(ctx context.Context, data model.AddComicTag, v *model.ComicTag) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic tag")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicTag(ctx, data, v)
}

func (svc Service) GetComicTagBySID(ctx context.Context, sid model.ComicTagSID) (*model.ComicTag, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var tagID any
	switch {
	case sid.TagID != nil:
		tagID = sid.TagID
	case sid.TagSID != nil:
		tagID = model.DBTagSIDToID(*sid.TagSID)
	}
	return svc.database.GetComicTag(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBTagGenericTagID:     tagID,
	})
}

func (svc Service) UpdateComicTagBySID(ctx context.Context, sid model.ComicTagSID, data model.SetComicTag, v *model.ComicTag) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic tag")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var tagID any
	switch {
	case sid.TagID != nil:
		tagID = sid.TagID
	case sid.TagSID != nil:
		tagID = model.DBTagSIDToID(*sid.TagSID)
	}
	return svc.database.UpdateComicTag(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBTagGenericTagID:     tagID,
	}, v)
}

func (svc Service) DeleteComicTagBySID(ctx context.Context, sid model.ComicTagSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic tag")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var tagID any
	switch {
	case sid.TagID != nil:
		tagID = sid.TagID
	case sid.TagSID != nil:
		tagID = model.DBTagSIDToID(*sid.TagSID)
	}
	return svc.database.DeleteComicTag(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBTagGenericTagID:     tagID,
	}, nil)
}

func (svc Service) ListComicTag(ctx context.Context, params model.ListParams) ([]*model.ComicTag, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicTagOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicTagOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicTagOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicTagPaginationMax {
			pagination.Limit = model.ComicTagPaginationMax
		}
	}

	return svc.database.ListComicTag(ctx, params)
}

func (svc Service) CountComicTag(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicTag(ctx, conds)
}

// Comic Relation

func (svc Service) AddComicRelationType(ctx context.Context, data model.AddComicRelationType, v *model.ComicRelationType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic relation type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicRelationType(ctx, data, v)
}

func (svc Service) GetComicRelationTypeByCode(ctx context.Context, code string) (*model.ComicRelationType, error) {
	return svc.database.GetComicRelationType(ctx, model.DBConditionalKV{
		Key:   model.DBComicRelationTypeCode,
		Value: code,
	})
}

func (svc Service) UpdateComicRelationTypeByCode(ctx context.Context, code string, data model.SetComicRelationType, v *model.ComicRelationType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic relation type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateComicRelationType(ctx, data, model.DBConditionalKV{
		Key:   model.DBComicRelationTypeCode,
		Value: code,
	}, v)
}

func (svc Service) DeleteComicRelationTypeByCode(ctx context.Context, code string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic relation type")
	}

	return svc.database.DeleteComicRelationType(ctx, model.DBConditionalKV{
		Key:   model.DBComicRelationTypeCode,
		Value: code,
	}, nil)
}

func (svc Service) ListComicRelationType(ctx context.Context, params model.ListParams) ([]*model.ComicRelationType, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicRelationTypeOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicRelationTypeOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicRelationTypeOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicRelationTypePaginationMax {
			pagination.Limit = model.ComicRelationTypePaginationMax
		}
	}

	return svc.database.ListComicRelationType(ctx, params)
}

func (svc Service) CountComicRelationType(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicRelationType(ctx, conds)
}

func (svc Service) AddComicRelation(ctx context.Context, data model.AddComicRelation, v *model.ComicRelation) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic relation")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicRelation(ctx, data, v)
}

func (svc Service) GetComicRelationBySID(ctx context.Context, sid model.ComicRelationSID) (*model.ComicRelation, error) {
	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBComicCodeToID(*sid.ParentCode)
	}
	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBComicRelationTypeCodeToID(*sid.TypeCode)
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBComicCodeToID(*sid.ChildCode)
	}
	return svc.database.GetComicRelation(ctx, map[string]any{
		model.DBComicRelationParentID: parentID,
		model.DBComicRelationTypeID:   typeID,
		model.DBComicRelationChildID:  childID,
	})
}

func (svc Service) UpdateComicRelationBySID(ctx context.Context, sid model.ComicRelationSID, data model.SetComicRelation, v *model.ComicRelation) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic relation")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBComicCodeToID(*sid.ParentCode)
	}
	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBComicRelationTypeCodeToID(*sid.TypeCode)
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBComicCodeToID(*sid.ChildCode)
	}
	return svc.database.UpdateComicRelation(ctx, data, map[string]any{
		model.DBComicRelationParentID: parentID,
		model.DBComicRelationTypeID:   typeID,
		model.DBComicRelationChildID:  childID,
	}, v)
}

func (svc Service) DeleteComicRelationBySID(ctx context.Context, sid model.ComicRelationSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic relation")
	}

	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBComicCodeToID(*sid.ParentCode)
	}
	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBComicRelationTypeCodeToID(*sid.TypeCode)
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBComicCodeToID(*sid.ChildCode)
	}
	return svc.database.DeleteComicRelation(ctx, map[string]any{
		model.DBComicRelationParentID: parentID,
		model.DBComicRelationTypeID:   typeID,
		model.DBComicRelationChildID:  childID,
	}, nil)
}

func (svc Service) ListComicRelation(ctx context.Context, params model.ListParams) ([]*model.ComicRelation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicRelationOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicRelationOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicRelationOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicRelationPaginationMax {
			pagination.Limit = model.ComicRelationPaginationMax
		}
	}

	return svc.database.ListComicRelation(ctx, params)
}

func (svc Service) CountComicRelation(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicRelation(ctx, conds)
}

//
// Comic Chapter
//

func (svc Service) AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic chapter")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicChapter(ctx, data, v)
}

func (svc Service) GetComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (*model.ComicChapter, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.database.GetComicChapter(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	})
}

func (svc Service) updateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic chapter")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateComicChapter(ctx, data, conds, v)
}

func (svc Service) UpdateComicChapterBySID(ctx context.Context, sid model.ComicChapterSID, data model.SetComicChapter, v *model.ComicChapter) error {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.updateComicChapter(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	}, v)
}

func (svc Service) deleteComicChapter(ctx context.Context, conds any, v *model.ComicChapter) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic chapter")
	}

	return svc.database.DeleteComicChapter(ctx, conds, v)
}

func (svc Service) DeleteComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) error {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.deleteComicChapter(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	}, nil)
}

func (svc Service) ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicChapterOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicChapterOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicChapterOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicChapterPaginationMax {
			pagination.Limit = model.ComicChapterPaginationMax
		}
	}

	return svc.database.ListComicChapter(ctx, params)
}

func (svc Service) CountComicChapter(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicChapter(ctx, conds)
}
