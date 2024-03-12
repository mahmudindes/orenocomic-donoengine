package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

const (
	NameErrCategoryTypeKey       = "category_type_code_key"
	NameErrCategoryKey           = "category_type_id_code_key"
	NameErrCategoryFKey          = "category_type_id_fkey"
	NameErrCategoryRelationPKey  = "category_relation_pkey"
	NameErrCategoryRelationFKey0 = "category_relation_parent_id_fkey"
	NameErrCategoryRelationFKey1 = "category_relation_child_id_fkey"
	NameErrCategoryRelationCheck = "category_relation_parent_id_child_id_check"
)

func (db Database) AddCategoryType(ctx context.Context, data model.AddCategoryType, v *model.CategoryType) error {
	if err := db.GenericAdd(ctx, model.DBCategoryType, map[string]any{
		model.DBCategoryTypeCode: data.Code,
		model.DBCategoryTypeName: data.Name,
	}, v); err != nil {
		return categoryTypeSetError(err)
	}
	return nil
}

func (db Database) GetCategoryType(ctx context.Context, conds any) (*model.CategoryType, error) {
	var result model.CategoryType
	if err := db.GenericGet(ctx, model.DBCategoryType, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateCategoryType(ctx context.Context, data model.SetCategoryType, conds any, v *model.CategoryType) error {
	data0 := map[string]any{}
	if data.Code != nil {
		data0[model.DBCategoryTypeCode] = data.Code
	}
	if data.Name != nil {
		data0[model.DBCategoryTypeName] = data.Name
	}
	if err := db.GenericUpdate(ctx, model.DBCategoryType, data0, conds, v); err != nil {
		return categoryTypeSetError(err)
	}
	return nil
}

func (db Database) DeleteCategoryType(ctx context.Context, conds any, v *model.CategoryType) error {
	return db.GenericDelete(ctx, model.DBCategoryType, conds, v)
}

func (db Database) ListCategoryType(ctx context.Context, params model.ListParams) ([]*model.CategoryType, error) {
	result := []*model.CategoryType{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBCategoryTypeCode})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.CategoryTypePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBCategoryType, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountCategoryType(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBCategoryType, conds)
}

func categoryTypeSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrCategoryTypeKey {
			return model.GenericError("same code already exists")
		}
	}
	return err
}

func (db Database) AddCategory(ctx context.Context, data model.AddCategory, v *model.Category) error {
	var typeID any
	switch {
	case data.TypeID != nil:
		typeID = data.TypeID
	case data.TypeCode != nil:
		typeID = model.DBCategoryTypeCodeToID(*data.TypeCode)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBCategoryTypeID: typeID,
		model.DBCategoryCode:   data.Code,
		model.DBCategoryName:   data.Name,
	})
	sql := "INSERT INTO " + model.DBCategory + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryTypeID + ", w." + model.DBCategoryCode
		sql += ", w." + model.DBCategoryName
		sql += ", l." + model.DBCategoryTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBCategoryType + " l"
		sql += " ON w." + model.DBCategoryTypeID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return categorySetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return categorySetError(err)
		}
	}
	return nil
}

func (db Database) GetCategory(ctx context.Context, conds any) (*model.Category, error) {
	var result model.Category
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBCategoryTypeID + ", w." + model.DBCategoryCode
	sql += ", w." + model.DBCategoryName
	sql += ", l." + model.DBCategoryTypeCode + " AS type_code"
	sql += " FROM " + model.DBCategory + " w JOIN " + model.DBCategoryType + " l"
	sql += " ON w." + model.DBCategoryTypeID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateCategory(ctx context.Context, data model.SetCategory, conds any, v *model.Category) error {
	txmode, data0 := false, map[string]any{}
	switch {
	case data.TypeID != nil:
		data0[model.DBCategoryTypeID] = data.TypeID
		if v == nil {
			v = new(model.Category)
		}
		txmode = true
	case data.TypeCode != nil:
		data0[model.DBCategoryTypeID] = model.DBCategoryTypeCodeToID(*data.TypeCode)
		if v == nil {
			v = new(model.Category)
		}
		txmode = true
	}
	if data.Code != nil {
		data0[model.DBCategoryCode] = data.Code
	}
	if data.Name != nil {
		data0[model.DBCategoryName] = data.Name
	}
	if txmode {
		ttx, err := db.ContextTransactionBegin(ctx)
		if err != nil {
			return err
		}
		ctx = ttx
		defer db.ContextTransactionRollback(context.WithoutCancel(ctx))
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBCategory + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryTypeID + ", w." + model.DBCategoryCode
		sql += ", w." + model.DBCategoryName
		sql += ", l." + model.DBCategoryTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBCategoryType + " l"
		sql += " ON w." + model.DBCategoryTypeID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return categorySetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return categorySetError(err)
		}
	}
	if data0[model.DBCategoryTypeID] != nil {
		typeChanged := false
		switch {
		case data.TypeID != nil && v.TypeID != *data.TypeID:
			typeChanged = true
		case data.TypeCode != nil && v.TypeCode != *data.TypeCode:
			typeChanged = true
		}
		if typeChanged {
			sql := "DELETE FROM " + model.DBCategoryRelation
			sql += " WHERE " + model.DBCategoryRelationParentID + " = " + utila.Utoa(v.ID)
			if err := db.Exec(ctx, sql); err != nil {
				return err
			}
		}
	}
	if txmode {
		if err := db.ContextTransactionCommit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) DeleteCategory(ctx context.Context, conds any, v *model.Category) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBCategory + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryTypeID + ", w." + model.DBCategoryCode
		sql += ", w." + model.DBCategoryName
		sql += ", l." + model.DBCategoryTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBCategoryType + " l"
		sql += " ON w." + model.DBCategoryTypeID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) ListCategory(ctx context.Context, params model.ListParams) ([]*model.Category, error) {
	result := []*model.Category{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBCategoryTypeID + ", w." + model.DBCategoryCode
	sql += ", w." + model.DBCategoryName
	sql += ", l." + model.DBCategoryTypeCode + " AS type_code"
	sql += " FROM " + model.DBCategory + " w JOIN " + model.DBCategoryType + " l"
	sql += " ON w." + model.DBCategoryTypeID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBCategoryCode})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.CategoryPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountCategory(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBCategory, conds)
}

func categorySetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign && errDatabase.Name == NameErrCategoryFKey {
			return model.GenericError("category type does not exist")
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrCategoryKey {
			return model.GenericError("same type id + code already exists")
		}
	}
	return err
}

func (db Database) AddCategoryRelation(ctx context.Context, data model.AddCategoryRelation, v *model.CategoryRelation) error {
	var parentID any
	switch {
	case data.ParentID != nil:
		parentID = data.ParentID
	case data.ParentCode != nil:
		parentID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.TypeID,
			TypeCode: data.TypeCode,
			Code:     *data.ParentCode,
		})
	}
	var childID any
	switch {
	case data.ChildID != nil:
		childID = data.ChildID
	case data.ChildCode != nil:
		childID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.TypeID,
			TypeCode: data.TypeCode,
			Code:     *data.ChildCode,
		})
	}
	if v == nil {
		v = new(model.CategoryRelation)
	}
	{
		ttx, err := db.ContextTransactionBegin(ctx)
		if err != nil {
			return err
		}
		ctx = ttx
		defer db.ContextTransactionRollback(context.WithoutCancel(ctx))
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBCategoryRelationParentID: parentID,
		model.DBCategoryRelationChildID:  childID,
	})
	sql := "INSERT INTO " + model.DBCategoryRelation + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryRelationParentID + ", w." + model.DBCategoryRelationChildID
		sql += ", l." + model.DBCategoryCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryRelationChildID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
	}
	{
		var parentLoop *bool
		args := []any{}
		cte := "WITH RECURSIVE childs AS ("
		cte += " SELECT " + model.DBCategoryRelationChildID + ", " + model.DBCategoryRelationParentID
		cte += "  FROM " + model.DBCategory
		cte += " UNION SELECT childs." + model.DBCategoryRelationChildID + ", parents." + model.DBCategoryRelationParentID
		cte += "  FROM " + model.DBCategory + " AS parents"
		cte += " JOIN childs ON childs." + model.DBCategoryRelationParentID + " = parents." + model.DBCategoryRelationChildID
		cte += ")"
		sql := cte + " SELECT (" + utila.Utoa(v.ParentID) + ", " + utila.Utoa(v.ChildID) + ") IN (SELECT * FROM childs)"
		if err := db.QueryOne(ctx, &parentLoop, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
		if parentLoop != nil && *parentLoop {
			return model.GenericError("category relation loop detected")
		}
	}
	return db.ContextTransactionCommit(ctx)
}

func (db Database) GetCategoryRelation(ctx context.Context, conds any) (*model.CategoryRelation, error) {
	var result model.CategoryRelation
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBCategoryRelationParentID + ", w." + model.DBCategoryRelationChildID
	sql += ", l." + model.DBCategoryCode + " AS child_code"
	sql += " FROM " + model.DBCategoryRelation + " w JOIN " + model.DBCategory + " l"
	sql += " ON w." + model.DBCategoryRelationChildID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateCategoryRelation(ctx context.Context, data model.SetCategoryRelation, conds any, v *model.CategoryRelation) error {
	data0 := map[string]any{}
	switch {
	case data.ParentID != nil:
		data0[model.DBCategoryRelationParentID] = data.ParentID
	case data.ParentCode != nil:
		data0[model.DBCategoryRelationParentID] = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.TypeID,
			TypeCode: data.TypeCode,
			Code:     *data.ParentCode,
		})
	}
	switch {
	case data.ChildID != nil:
		data0[model.DBCategoryRelationChildID] = data.ChildID
	case data.ChildCode != nil:
		data0[model.DBCategoryRelationChildID] = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.TypeID,
			TypeCode: data.TypeCode,
			Code:     *data.ChildCode,
		})
	}
	if v == nil {
		v = new(model.CategoryRelation)
	}
	{
		ttx, err := db.ContextTransactionBegin(ctx)
		if err != nil {
			return err
		}
		ctx = ttx
		defer db.ContextTransactionRollback(context.WithoutCancel(ctx))
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicRelation + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryRelationParentID + ", w." + model.DBCategoryRelationChildID
		sql += ", l." + model.DBCategoryCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryRelationChildID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
	}
	{
		var parentLoop *bool
		args := []any{}
		cte := "WITH RECURSIVE childs AS ("
		cte += " SELECT " + model.DBCategoryRelationChildID + ", " + model.DBCategoryRelationParentID
		cte += "  FROM " + model.DBCategory
		cte += " UNION SELECT childs." + model.DBCategoryRelationChildID + ", parents." + model.DBCategoryRelationParentID
		cte += "  FROM " + model.DBCategory + " AS parents"
		cte += " JOIN childs ON childs." + model.DBCategoryRelationParentID + " = parents." + model.DBCategoryRelationChildID
		cte += ")"
		sql := cte + " SELECT (" + utila.Utoa(v.ParentID) + ", " + utila.Utoa(v.ChildID) + ") IN (SELECT * FROM childs)"
		if err := db.QueryOne(ctx, &parentLoop, sql, args...); err != nil {
			return categoryRelationSetError(err)
		}
		if parentLoop != nil && *parentLoop {
			return model.GenericError("category relation loop detected")
		}
	}
	return db.ContextTransactionCommit(ctx)
}

func (db Database) DeleteCategoryRelation(ctx context.Context, conds any, v *model.CategoryRelation) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBCategoryRelation + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBCategoryRelationParentID + ", w." + model.DBCategoryRelationChildID
		sql += ", l." + model.DBCategoryCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryRelationChildID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) ListCategoryRelation(ctx context.Context, params model.ListParams) ([]*model.CategoryRelation, error) {
	result := []*model.CategoryRelation{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBCategoryRelationParentID + ", w." + model.DBCategoryRelationChildID
	sql += ", l." + model.DBCategoryCode + " AS child_code"
	sql += " FROM " + model.DBCategoryRelation + " w JOIN " + model.DBCategory + " l"
	sql += " ON w." + model.DBCategoryRelationChildID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBCategoryRelationChildID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.CategoryRelationPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountCategoryRelation(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBCategoryRelation, conds)
}

func categoryRelationSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrCategoryRelationFKey0:
				return model.GenericError("parent category does not exist")
			case NameErrCategoryRelationFKey1:
				return model.GenericError("child category does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrCategoryRelationPKey {
			return model.GenericError("same child id already exists")
		}
		if errDatabase.Code == CodeErrValidation && errDatabase.Name == NameErrCategoryRelationCheck {
			return model.GenericError("parent category and child category cannot be same")
		}
	}
	return err
}
