package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

func (svc Service) AddCategoryType(ctx context.Context, data model.AddCategoryType, v *model.CategoryType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add category type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddCategoryType(ctx, data, v)
}

func (svc Service) GetCategoryTypeByCode(ctx context.Context, code string) (*model.CategoryType, error) {
	return svc.database.GetCategoryType(ctx, model.DBConditionalKV{
		Key:   model.DBCategoryTypeCode,
		Value: code,
	})
}

func (svc Service) UpdateCategoryTypeByCode(ctx context.Context, code string, data model.SetCategoryType, v *model.CategoryType) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update category type")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateCategoryType(ctx, data, model.DBConditionalKV{
		Key:   model.DBCategoryTypeCode,
		Value: code,
	}, v)
}

func (svc Service) DeleteCategoryTypeByCode(ctx context.Context, code string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete category type")
	}

	return svc.database.DeleteCategoryType(ctx, model.DBConditionalKV{
		Key:   model.DBCategoryTypeCode,
		Value: code,
	}, nil)
}

func (svc Service) ListCategoryType(ctx context.Context, params model.ListParams) ([]*model.CategoryType, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.CategoryTypeOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.CategoryTypeOrderBysMax {
		params.OrderBys = params.OrderBys[:model.CategoryTypeOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.CategoryTypePaginationMax {
			pagination.Limit = model.CategoryTypePaginationMax
		}
	}

	return svc.database.ListCategoryType(ctx, params)
}

func (svc Service) CountCategoryType(ctx context.Context, conds any) (int, error) {
	return svc.database.CountCategoryType(ctx, conds)
}

func (svc Service) AddCategory(ctx context.Context, data model.AddCategory, v *model.Category) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add category")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.Relations = []*model.CategoryRelation{}
	}

	return svc.database.AddCategory(ctx, data, v)
}

func (svc Service) GetCategoryBySID(ctx context.Context, sid model.CategorySID) (*model.Category, error) {
	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBCategoryTypeCodeToID(*sid.TypeCode)
	}
	result, err := svc.database.GetCategory(ctx, map[string]any{
		model.DBCategoryTypeID: typeID,
		model.DBCategoryCode:   sid.Code,
	})
	if err != nil {
		return nil, err
	}

	relations, err := svc.ListCategoryRelation(ctx, model.ListParams{
		Conditions: model.DBConditionalKV{Key: model.DBCategoryRelationParentID, Value: result.ID},
		Pagination: &model.Pagination{},
	})
	if err != nil {
		return nil, err
	}
	result.Relations = relations

	return result, nil
}

func (svc Service) UpdateCategoryBySID(ctx context.Context, sid model.CategorySID, data model.SetCategory, v *model.Category) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update category")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBCategoryTypeCodeToID(*sid.TypeCode)
	}
	if err := svc.database.UpdateCategory(ctx, data, map[string]any{
		model.DBCategoryTypeID: typeID,
		model.DBCategoryCode:   sid.Code,
	}, v); err != nil {
		return err
	}

	if v != nil {
		relations, err := svc.ListCategoryRelation(ctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBCategoryRelationParentID, Value: v.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}
		v.Relations = relations
	}

	return nil
}

func (svc Service) DeleteCategoryBySID(ctx context.Context, sid model.CategorySID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete category")
	}

	var typeID any
	switch {
	case sid.TypeID != nil:
		typeID = sid.TypeID
	case sid.TypeCode != nil:
		typeID = model.DBCategoryTypeCodeToID(*sid.TypeCode)
	}
	return svc.database.DeleteCategory(ctx, map[string]any{
		model.DBCategoryTypeID: typeID,
		model.DBCategoryCode:   sid.Code,
	}, nil)
}

func (svc Service) ListCategory(ctx context.Context, params model.ListParams) ([]*model.Category, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.CategoryOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.CategoryOrderBysMax {
		params.OrderBys = params.OrderBys[:model.CategoryOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.CategoryPaginationMax {
			pagination.Limit = model.CategoryPaginationMax
		}
	}

	result, err := svc.database.ListCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBCategoryRelationParentID,
				Value: r.ID,
			})
		}
		relations, err := svc.database.ListCategoryRelation(ctx, model.ListParams{
			Conditions: conds,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		for _, r := range result {
			r.Relations = []*model.CategoryRelation{}
		}
		for _, relation := range relations {
			for _, r := range result {
				if r.ID == relation.ParentID {
					r.Relations = append(r.Relations, relation)
				}
			}
		}
	}

	return result, nil
}

func (svc Service) CountCategory(ctx context.Context, conds any) (int, error) {
	return svc.database.CountCategory(ctx, conds)
}

func (svc Service) AddCategoryRelation(ctx context.Context, data model.AddCategoryRelation, v *model.CategoryRelation) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add category relation")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddCategoryRelation(ctx, data, v)
}

func (svc Service) GetCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID) (*model.CategoryRelation, error) {
	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ParentCode,
		})
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ChildCode,
		})
	}
	return svc.database.GetCategoryRelation(ctx, map[string]any{
		model.DBCategoryRelationParentID: parentID,
		model.DBCategoryRelationChildID:  childID,
	})
}

func (svc Service) UpdateCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID, data model.SetCategoryRelation, v *model.CategoryRelation) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update category relation")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ChildCode,
		})
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ChildCode,
		})
	}
	return svc.database.UpdateCategoryRelation(ctx, data, map[string]any{
		model.DBCategoryRelationParentID: parentID,
		model.DBCategoryRelationChildID:  childID,
	}, v)
}

func (svc Service) DeleteCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete category relation")
	}

	var parentID any
	switch {
	case sid.ParentID != nil:
		parentID = sid.ParentID
	case sid.ParentCode != nil:
		parentID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ChildCode,
		})
	}
	var childID any
	switch {
	case sid.ChildID != nil:
		childID = sid.ChildID
	case sid.ChildCode != nil:
		childID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   sid.TypeID,
			TypeCode: sid.TypeCode,
			Code:     *sid.ChildCode,
		})
	}
	return svc.database.DeleteCategoryRelation(ctx, map[string]any{
		model.DBCategoryRelationParentID: parentID,
		model.DBCategoryRelationChildID:  childID,
	}, nil)
}

func (svc Service) ListCategoryRelation(ctx context.Context, params model.ListParams) ([]*model.CategoryRelation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.CategoryRelationOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.CategoryRelationOrderBysMax {
		params.OrderBys = params.OrderBys[:model.CategoryRelationOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.CategoryRelationPaginationMax {
			pagination.Limit = model.CategoryRelationPaginationMax
		}
	}

	return svc.database.ListCategoryRelation(ctx, params)
}

func (svc Service) CountCategoryRelation(ctx context.Context, conds any) (int, error) {
	return svc.database.CountCategoryRelation(ctx, conds)
}
