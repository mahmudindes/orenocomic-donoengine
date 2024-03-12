package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

func (svc Service) AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add website")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddWebsite(ctx, data, v)
}

func (svc Service) GetWebsiteByDomain(ctx context.Context, domain string) (*model.Website, error) {
	return svc.database.GetWebsite(ctx, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	})
}

func (svc Service) UpdateWebsiteByDomain(ctx context.Context, domain string, data model.SetWebsite, v *model.Website) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update website")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateWebsite(ctx, data, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	}, v)
}

func (svc Service) DeleteWebsiteByDomain(ctx context.Context, domain string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete website")
	}

	return svc.database.DeleteWebsite(ctx, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	}, nil)
}

func (svc Service) ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.WebsiteOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.WebsiteOrderBysMax {
		params.OrderBys = params.OrderBys[:model.WebsiteOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.WebsitePaginationMax {
			pagination.Limit = model.WebsitePaginationMax
		}
	}

	return svc.database.ListWebsite(ctx, params)
}

func (svc Service) CountWebsite(ctx context.Context, conds any) (int, error) {
	return svc.database.CountWebsite(ctx, conds)
}
