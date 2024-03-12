package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

const (
	NameErrTagTypeKey = "tag_type_code_key"
	NameErrTagKey     = "tag_type_id_code_key"
	NameErrTagFKey    = "tag_type_id_fkey"
)

func (db Database) AddTagType(ctx context.Context, data model.AddTagType, v *model.TagType) error {
	if err := db.GenericAdd(ctx, model.DBTagType, map[string]any{
		model.DBTagTypeCode: data.Code,
		model.DBTagTypeName: data.Name,
	}, v); err != nil {
		return tagTypeSetError(err)
	}
	return nil
}

func (db Database) GetTagType(ctx context.Context, conds any) (*model.TagType, error) {
	var result model.TagType
	if err := db.GenericGet(ctx, model.DBTagType, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateTagType(ctx context.Context, data model.SetTagType, conds any, v *model.TagType) error {
	data0 := map[string]any{}
	if data.Code != nil {
		data0[model.DBTagTypeCode] = data.Code
	}
	if data.Name != nil {
		data0[model.DBTagTypeName] = data.Name
	}
	if err := db.GenericUpdate(ctx, model.DBTagType, data0, conds, v); err != nil {
		return tagTypeSetError(err)
	}
	return nil
}

func (db Database) DeleteTagType(ctx context.Context, conds any, v *model.TagType) error {
	return db.GenericDelete(ctx, model.DBTagType, conds, v)
}

func (db Database) ListTagType(ctx context.Context, params model.ListParams) ([]*model.TagType, error) {
	result := []*model.TagType{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBTagTypeCode})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.TagTypePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBTagType, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountTagType(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBTagType, conds)
}

func tagTypeSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrTagTypeKey {
			return model.GenericError("same code already exists")
		}
	}
	return err
}

func (db Database) AddTag(ctx context.Context, data model.AddTag, v *model.Tag) error {
	var typeID any
	switch {
	case data.TypeID != nil:
		typeID = data.TypeID
	case data.TypeCode != nil:
		typeID = model.DBTagTypeCodeToID(*data.TypeCode)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBTagTypeID: typeID,
		model.DBTagCode:   data.Code,
		model.DBTagName:   data.Name,
	})
	sql := "INSERT INTO " + model.DBTag + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBTagTypeID + ", w." + model.DBTagCode + ", w." + model.DBTagName
		sql += ", l." + model.DBTagTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBTagType + " l"
		sql += " ON w." + model.DBTagTypeID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return tagSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return tagSetError(err)
		}
	}
	return nil
}

func (db Database) GetTag(ctx context.Context, conds any) (*model.Tag, error) {
	var result model.Tag
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBTagTypeID + ", w." + model.DBTagCode + ", w." + model.DBTagName
	sql += ", l." + model.DBTagTypeCode + " AS type_code"
	sql += " FROM " + model.DBTag + " w JOIN " + model.DBTagType + " l"
	sql += " ON w." + model.DBTagTypeID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateTag(ctx context.Context, data model.SetTag, conds any, v *model.Tag) error {
	data0 := map[string]any{}
	switch {
	case data.TypeID != nil:
		data0[model.DBCategoryTypeID] = data.TypeID
	case data.TypeCode != nil:
		data0[model.DBCategoryTypeID] = model.DBTagTypeCodeToID(*data.TypeCode)
	}
	if data.Code != nil {
		data0[model.DBCategoryCode] = data.Code
	}
	if data.Name != nil {
		data0[model.DBCategoryName] = data.Name
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBTag + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBTagTypeID + ", w." + model.DBTagCode + ", w." + model.DBTagName
		sql += ", l." + model.DBTagTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBTagType + " l"
		sql += " ON w." + model.DBTagTypeID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return tagSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return tagSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteTag(ctx context.Context, conds any, v *model.Tag) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBTag + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBTagTypeID + ", w." + model.DBTagCode + ", w." + model.DBTagName
		sql += ", l." + model.DBTagTypeCode + " AS type_code"
		sql += " FROM data w JOIN " + model.DBTagType + " l"
		sql += " ON w." + model.DBTagTypeID + " = l." + model.DBGenericID
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

func (db Database) ListTag(ctx context.Context, params model.ListParams) ([]*model.Tag, error) {
	result := []*model.Tag{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBTagTypeID + ", w." + model.DBTagCode + ", w." + model.DBTagName
	sql += ", l." + model.DBTagTypeCode + " AS type_code"
	sql += " FROM " + model.DBTag + " w JOIN " + model.DBTagType + " l"
	sql += " ON w." + model.DBTagTypeID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBTagCode})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.TagPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountTag(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBTag, conds)
}

func tagSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign && errDatabase.Name == NameErrTagFKey {
			return model.GenericError("tag type does not exist")
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrTagKey {
			return model.GenericError("same type id + code already exists")
		}
	}
	return err
}
