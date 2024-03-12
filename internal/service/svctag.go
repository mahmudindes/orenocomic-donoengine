package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

func (svc Service) AddTagType(ctx context.Context, data model.AddTagType, v *model.TagType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add tag type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddTagType(ctx, data, v)
}

func (svc Service) GetTagTypeByCode(ctx context.Context, code string) (*model.TagType, error) {
	return svc.database.GetTagType(ctx, model.DBConditionalKV{
		Key:   model.DBTagTypeCode,
		Value: code,
	})
}

func (svc Service) UpdateTagTypeByCode(ctx context.Context, code string, data model.SetTagType, v *model.TagType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update tag type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateTagType(ctx, data, model.DBConditionalKV{
		Key:   model.DBTagTypeCode,
		Value: code,
	}, v)
}

func (svc Service) DeleteTagTypeByCode(ctx context.Context, code string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete tag type")
	}

	return svc.database.DeleteTagType(ctx, model.DBConditionalKV{
		Key:   model.DBTagTypeCode,
		Value: code,
	}, nil)
}

func (svc Service) ListTagType(ctx context.Context, params model.ListParams) ([]*model.TagType, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.TagTypeOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.TagTypeOrderBysMax {
		params.OrderBys = params.OrderBys[:model.TagTypeOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.TagTypePaginationMax {
			pagination.Limit = model.TagTypePaginationMax
		}
	}

	return svc.database.ListTagType(ctx, params)
}

func (svc Service) CountTagType(ctx context.Context, conds any) (int, error) {
	return svc.database.CountTagType(ctx, conds)
}

func (svc Service) AddTag(ctx context.Context, data model.AddTag, v *model.Tag) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add tag")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddTag(ctx, data, v)
}

func (svc Service) GetTagBySID(ctx context.Context, sid model.TagSID) (*model.Tag, error) {
	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBTagTypeCodeToID(*sid.TypeCode)
	}
	return svc.database.GetTag(ctx, map[string]any{
		model.DBTagTypeID: typeID,
		model.DBTagCode:   sid.Code,
	})
}

func (svc Service) UpdateTagBySID(ctx context.Context, sid model.TagSID, data model.SetTag, v *model.Tag) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update tag")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBTagTypeCodeToID(*sid.TypeCode)
	}
	return svc.database.UpdateTag(ctx, data, map[string]any{
		model.DBTagTypeID: typeID,
		model.DBTagCode:   sid.Code,
	}, v)
}

func (svc Service) DeleteTagBySID(ctx context.Context, sid model.TagSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete tag")
	}

	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBTagTypeCodeToID(*sid.TypeCode)
	}
	return svc.database.DeleteTag(ctx, map[string]any{
		model.DBTagTypeID: typeID,
		model.DBTagCode:   sid.Code,
	}, nil)
}

func (svc Service) ListTag(ctx context.Context, params model.ListParams) ([]*model.Tag, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.TagOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.TagOrderBysMax {
		params.OrderBys = params.OrderBys[:model.TagOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.TagPaginationMax {
			pagination.Limit = model.TagPaginationMax
		}
	}

	return svc.database.ListTag(ctx, params)
}

func (svc Service) CountTag(ctx context.Context, conds any) (int, error) {
	return svc.database.CountTag(ctx, conds)
}
