package database

import (
	"context"
	"errors"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

const (
	NameErrWebsiteKey = "website_domain_key"
)

func (db Database) AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error {
	if err := db.GenericAdd(ctx, model.DBWebsite, map[string]any{
		model.DBWebsiteDomain: data.Domain,
		model.DBLanguageName:  data.Name,
	}, v); err != nil {
		return websiteSetError(err)
	}
	return nil
}

func (db Database) GetWebsite(ctx context.Context, conds any) (*model.Website, error) {
	var result model.Website
	if err := db.GenericGet(ctx, model.DBWebsite, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateWebsite(ctx context.Context, data model.SetWebsite, conds any, v *model.Website) error {
	data0 := map[string]any{}
	if data.Domain != nil {
		data0[model.DBWebsiteDomain] = data.Domain
	}
	if data.Name != nil {
		data0[model.DBWebsiteName] = data.Name
	}
	if err := db.GenericUpdate(ctx, model.DBWebsite, data0, conds, v); err != nil {
		return websiteSetError(err)
	}
	return nil
}

func (db Database) DeleteWebsite(ctx context.Context, conds any, v *model.Website) error {
	return db.GenericDelete(ctx, model.DBWebsite, conds, v)
}

func (db Database) ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error) {
	result := []*model.Website{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBWebsiteDomain})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.WebsitePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBWebsite, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountWebsite(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBWebsite, conds)
}

func websiteSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrWebsiteKey {
			return model.GenericError("same domain already exists")
		}
	}
	return err
}
