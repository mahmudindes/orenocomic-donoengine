package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

const (
	NameErrComicKey  = "comic_code_key"
	NameErrComicFKey = "comic_language_id_fkey"
)

func (db Database) AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error {
	var code any
	switch {
	case data.Code != nil:
		code = data.Code
	default:
		code = utila.RandomString(utila.RandomStringGeneral, model.ComicCodeLength)
	}
	var languageID any
	switch {
	case data.LanguageID != nil:
		languageID = data.LanguageID
	case data.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicCode:                 code,
		model.DBLanguageGenericLanguageID: languageID,
		model.DBComicPublishedFrom:        data.PublishedFrom,
		model.DBComicPublishedTo:          data.PublishedTo,
		model.DBComicTotalChapter:         data.TotalChapter,
		model.DBComicTotalVolume:          data.TotalVolume,
		model.DBComicNSFW:                 data.NSFW,
		model.DBComicNSFL:                 data.NSFL,
		model.DBComicAdditionals:          data.Additionals,
	})
	sql := "INSERT INTO " + model.DBComic + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicCode + ", w." + model.DBComicPublishedFrom
		sql += ", w." + model.DBComicPublishedTo + ", w." + model.DBComicTotalChapter
		sql += ", w." + model.DBComicTotalVolume + ", w." + model.DBComicNSFW
		sql += ", w." + model.DBComicNSFL + ", w." + model.DBLanguageGenericLanguageID
		sql += ", w." + model.DBComicAdditionals
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w LEFT JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicSetError(err)
		}
	}
	return nil
}

func (db Database) GetComic(ctx context.Context, conds any) (*model.Comic, error) {
	var result model.Comic
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicCode + ", w." + model.DBComicPublishedFrom
	sql += ", w." + model.DBComicPublishedTo + ", w." + model.DBComicTotalChapter
	sql += ", w." + model.DBComicTotalVolume + ", w." + model.DBComicNSFW
	sql += ", w." + model.DBComicNSFL + ", w." + model.DBLanguageGenericLanguageID
	sql += ", w." + model.DBComicAdditionals
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComic + " w LEFT JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComic(ctx context.Context, data model.SetComic, conds any, v *model.Comic) error {
	data0 := map[string]any{}
	if data.Code != nil {
		data0[model.DBComicCode] = data.Code
	}
	switch {
	case data.LanguageID != nil:
		data0[model.DBLanguageGenericLanguageID] = data.LanguageID
	case data.LanguageIETF != nil:
		data0[model.DBLanguageGenericLanguageID] = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	if data.PublishedFrom != nil {
		data0[model.DBComicPublishedFrom] = data.PublishedFrom
	}
	if data.PublishedTo != nil {
		data0[model.DBComicPublishedTo] = data.PublishedTo
	}
	if data.TotalChapter != nil {
		data0[model.DBComicTotalChapter] = data.TotalChapter
	}
	if data.TotalVolume != nil {
		data0[model.DBComicTotalVolume] = data.TotalVolume
	}
	if data.NSFW != nil {
		data0[model.DBComicNSFW] = data.NSFW
	}
	if data.NSFL != nil {
		data0[model.DBComicNSFL] = data.NSFL
	}
	if data.Additionals != nil {
		data0[model.DBComicAdditionals] = data.Additionals
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComic + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicCode + ", w." + model.DBComicPublishedFrom
		sql += ", w." + model.DBComicPublishedTo + ", w." + model.DBComicTotalChapter
		sql += ", w." + model.DBComicTotalVolume + ", w." + model.DBComicNSFW
		sql += ", w." + model.DBComicNSFL + ", w." + model.DBLanguageGenericLanguageID
		sql += ", w." + model.DBComicAdditionals
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w LEFT JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComic(ctx context.Context, conds any, v *model.Comic) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComic + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicCode + ", w." + model.DBComicPublishedFrom
		sql += ", w." + model.DBComicPublishedTo + ", w." + model.DBComicTotalChapter
		sql += ", w." + model.DBComicTotalVolume + ", w." + model.DBComicNSFW
		sql += ", w." + model.DBComicNSFL + ", w." + model.DBLanguageGenericLanguageID
		sql += ", w." + model.DBComicAdditionals
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w LEFT JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
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

func (db Database) ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error) {
	result := []*model.Comic{}
	args := []any{}
	ccnd := map[string]model.DBCrossConditional{}
	switch cond := params.Conditions.(type) {
	case model.DBCrossConditional:
		switch cond.Table {
		case model.DBComicExternal:
			ccnd["ce"] = cond
		}
	case []any:
		for _, cond := range cond {
			switch cond := cond.(type) {
			case model.DBCrossConditional:
				switch cond.Table {
				case model.DBComicExternal:
					ccnd["ce"] = cond
				}
			}
		}
	}
	sql := "SELECT * FROM ("
	if len(ccnd) > 0 {
		cte := ""
		for key, val := range ccnd {
			switch val.Table {
			case model.DBComicExternal:
				if cte != "" {
					cte += ", "
				}
				cte += key + "cte AS("
				cte += "SELECT * FROM (SELECT a." + model.DBGenericID
				cte += ", a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
				cte += ", a." + model.DBComicGenericComicID + ", a." + model.DBComicGenericRID
				cte += ", a." + model.DBWebsiteGenericWebsiteID + ", a." + model.DBComicExternalRelativeURL
				cte += ", a." + model.DBComicExternalOfficial
				cte += ", b." + model.DBWebsiteDomain + " AS website_domain"
				cte += " FROM " + model.DBComicExternal + " a JOIN " + model.DBWebsite + " b"
				cte += " ON a." + model.DBWebsiteGenericWebsiteID + " = b." + model.DBGenericID
				cte += ")"
				if cond := SetWhere(val.Conditions, &args); cond != "" {
					cte += " WHERE " + cond
				}
				cte += ")"
			}
		}
		if cte != "" {
			sql += "WITH " + cte + " "
		}
		sql += "SELECT DISTINCT ON (a." + model.DBGenericID + ")"
		sql += " a." + model.DBGenericID
	} else {
		sql += "SELECT a." + model.DBGenericID
	}
	sql += ", a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
	sql += ", a." + model.DBComicCode + ", a." + model.DBComicPublishedFrom
	sql += ", a." + model.DBComicPublishedTo + ", a." + model.DBComicTotalChapter
	sql += ", a." + model.DBComicTotalVolume + ", a." + model.DBComicNSFW
	sql += ", a." + model.DBComicNSFL + ", a." + model.DBLanguageGenericLanguageID
	sql += ", a." + model.DBComicAdditionals
	sql += ", b." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComic + " a LEFT JOIN " + model.DBLanguage + " b"
	sql += " ON a." + model.DBLanguageGenericLanguageID + " = b." + model.DBGenericID
	for key, val := range ccnd {
		switch val.Table {
		case model.DBComicExternal:
			sql += " LEFT JOIN " + key + "cte " + key
			sql += " ON a." + model.DBGenericID + " = " + key + "." + model.DBComicGenericComicID
		}
	}
	if len(ccnd) > 0 {
		whr := ""
		for key, val := range ccnd {
			switch val.Table {
			case model.DBComicExternal:
				if whr != "" {
					whr += " AND"
				}
				whr += " " + key + "." + model.DBComicGenericComicID + " IS NOT NULL"
			}
		}
		if whr != "" {
			sql += " WHERE" + whr
		}
	}
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicCode})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComic(ctx context.Context, conds any) (int, error) {
	var dst int
	args := []any{}
	ccnd := map[string]model.DBCrossConditional{}
	switch cond := conds.(type) {
	case model.DBCrossConditional:
		switch cond.Table {
		case model.DBComicExternal:
			ccnd["ce"] = cond
		}
	case []any:
		for _, cond := range cond {
			switch cond := cond.(type) {
			case model.DBCrossConditional:
				switch cond.Table {
				case model.DBComicExternal:
					ccnd["ce"] = cond
				}
			}
		}
	}
	sql := ""
	if len(ccnd) > 0 {
		cte := ""
		for key, val := range ccnd {
			switch val.Table {
			case model.DBComicExternal:
				if cte != "" {
					cte += ", "
				}
				cte += key + "cte AS("
				cte += "SELECT * FROM (SELECT a." + model.DBGenericID
				cte += ", a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
				cte += ", a." + model.DBComicGenericComicID + ", a." + model.DBComicGenericRID
				cte += ", a." + model.DBWebsiteGenericWebsiteID + ", a." + model.DBComicExternalRelativeURL
				cte += ", a." + model.DBComicExternalOfficial
				cte += ", b." + model.DBWebsiteDomain + " AS website_domain"
				cte += " FROM " + model.DBComicExternal + " a JOIN " + model.DBWebsite + " b"
				cte += " ON a." + model.DBWebsiteGenericWebsiteID + " = b." + model.DBGenericID
				cte += ")"
				if cond := SetWhere(val.Conditions, &args); cond != "" {
					cte += " WHERE " + cond
				}
				cte += ")"
			}
		}
		if cte != "" {
			sql += "WITH " + cte + " "
		}
		sql += "SELECT COUNT(DISTINCT a." + model.DBGenericID + ")"
		sql += " FROM " + model.DBComic + " a"
		for key, val := range ccnd {
			switch val.Table {
			case model.DBComicExternal:
				sql += " LEFT JOIN " + key + "cte " + key
				sql += " ON a." + model.DBGenericID + " = " + key + "." + model.DBComicGenericComicID
			}
		}
		whr := SetWhere(conds, &args)
		for key, val := range ccnd {
			switch val.Table {
			case model.DBComicExternal:
				if whr != "" {
					whr += " AND"
				}
				whr += " " + key + "." + model.DBComicGenericComicID + " IS NOT NULL"
			}
		}
		if whr != "" {
			sql += " WHERE" + whr
		}
	} else {
		sql += "SELECT COUNT(*) FROM " + model.DBComic
		if cond := SetWhere(conds, &args); cond != "" {
			sql += " WHERE " + cond
		}
	}
	if err := db.QueryOne(ctx, &dst, sql, args...); err != nil {
		return -1, err
	}
	return dst, nil
}

func (db Database) ExistsComic(ctx context.Context, conds any) (bool, error) {
	return db.GenericExists(ctx, model.DBComic, conds)
}

func comicSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign && errDatabase.Name == NameErrComicFKey {
			return model.GenericError("language does not exist")
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicKey {
			return model.GenericError("same code already exists")
		}
	}
	return err
}

const (
	NameErrComicTitleKey0       = "comic_title_comic_id_rid_key"
	NameErrComicTitleKey1       = "comic_title_comic_id_title_key"
	NameErrComicTitleFKey0      = "comic_title_comic_id_fkey"
	NameErrComicTitleFKey1      = "comic_title_language_id_fkey"
	NameErrComicCoverKey0       = "comic_cover_comic_id_rid_key"
	NameErrComicCoverKey1       = "comic_cover_comic_id_website_id_relative_url_key"
	NameErrComicCoverFKey0      = "comic_cover_comic_id_fkey"
	NameErrComicCoverFKey1      = "comic_cover_website_id_fkey"
	NameErrComicSynopsisKey0    = "comic_synopsis_comic_id_rid_key"
	NameErrComicSynopsisKey1    = "comic_synopsis_comic_id_synopsis_key"
	NameErrComicSynopsisFKey0   = "comic_synopsis_comic_id_fkey"
	NameErrComicSynopsisFKey1   = "comic_synopsis_language_id_fkey"
	NameErrComicExternalKey0    = "comic_external_comic_id_rid_key"
	NameErrComicExternalKey1    = "comic_external_comic_id_website_id_relative_url_key"
	NameErrComicExternalFKey0   = "comic_external_comic_id_fkey"
	NameErrComicExternalFKey1   = "comic_external_website_id_fkey"
	NameErrComicCategoryPKey    = "comic_category_pkey"
	NameErrComicCategoryFKey0   = "comic_category_comic_id_fkey"
	NameErrComicCategoryFKey1   = "comic_category_category_id_fkey"
	NameErrComicTagPKey         = "comic_tag_pkey"
	NameErrComicTagFKey0        = "comic_tag_comic_id_fkey"
	NameErrComicTagFKey1        = "comic_tag_tag_id_fkey"
	NameErrComicRelationTypeKey = "comic_relation_type_code_key"
	NameErrComicRelationPKey    = "comic_relation_pkey"
	NameErrComicRelationFKey0   = "comic_relation_type_id_fkey"
	NameErrComicRelationFKey1   = "comic_relation_parent_id_fkey"
	NameErrComicRelationFKey2   = "comic_relation_child_id_fkey"
	NameErrComicRelationCheck   = "comic_relation_parent_id_child_id_check"
)

func (db Database) AddComicTitle(ctx context.Context, data model.AddComicTitle, v *model.ComicTitle) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var rid any
	switch {
	case data.RID != nil:
		rid = data.RID
	default:
		rid = utila.RandomString(utila.RandomStringGeneral, model.ComicGenericRIDLength)
	}
	var languageID any
	switch {
	case data.LanguageID != nil:
		languageID = data.LanguageID
	case data.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBComicGenericRID:           rid,
		model.DBLanguageGenericLanguageID: languageID,
		model.DBComicTitleTitle:           data.Title,
		model.DBComicTitleSynonym:         data.Synonym,
		model.DBComicTitleRomanized:       data.Romanized,
	})
	sql := "INSERT INTO " + model.DBComicTitle + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicTitleTitle
		sql += ", w." + model.DBComicTitleSynonym + ", w." + model.DBComicTitleRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicTitleSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicTitleSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicTitle(ctx context.Context, conds any) (*model.ComicTitle, error) {
	var result model.ComicTitle
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicTitleTitle
	sql += ", w." + model.DBComicTitleSynonym + ", w." + model.DBComicTitleRomanized
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComicTitle + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicTitle(ctx context.Context, data model.SetComicTitle, conds any, v *model.ComicTitle) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.RID != nil {
		data0[model.DBComicGenericRID] = data.RID
	}
	switch {
	case data.LanguageID != nil:
		data0[model.DBLanguageGenericLanguageID] = data.LanguageID
	case data.LanguageIETF != nil:
		data0[model.DBLanguageGenericLanguageID] = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	if data.Title != nil {
		data0[model.DBComicTitleTitle] = data.Title
	}
	if data.Synonym != nil {
		data0[model.DBComicTitleSynonym] = data.Synonym
	}
	if data.Romanized != nil {
		data0[model.DBComicTitleRomanized] = data.Romanized
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicTitle + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicTitleTitle
		sql += ", w." + model.DBComicTitleSynonym + ", w." + model.DBComicTitleRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicTitleSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicTitleSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicTitle(ctx context.Context, conds any, v *model.ComicTitle) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicTitle + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicTitleTitle
		sql += ", w." + model.DBComicTitleSynonym + ", w." + model.DBComicTitleRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
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

func (db Database) ListComicTitle(ctx context.Context, params model.ListParams) ([]*model.ComicTitle, error) {
	result := []*model.ComicTitle{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicTitleTitle
	sql += ", w." + model.DBComicTitleSynonym + ", w." + model.DBComicTitleRomanized
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComicTitle + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicGenericRID})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicTitlePaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicTitle(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicTitle, conds)
}

func comicTitleSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicTitleFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicTitleFKey1:
				return model.GenericError("language does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists {
			switch errDatabase.Name {
			case NameErrComicTitleKey0:
				return model.GenericError("same comic id + rid already exists")
			case NameErrComicTitleKey1:
				return model.GenericError("same comic id + title already exists")
			}
		}
	}
	return err
}

func (db Database) AddComicCover(ctx context.Context, data model.AddComicCover, v *model.ComicCover) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var rid any
	switch {
	case data.RID != nil:
		rid = data.RID
	default:
		rid = utila.RandomString(utila.RandomStringGeneral, model.ComicGenericRIDLength)
	}
	var websiteID any
	switch {
	case data.WebsiteID != nil:
		websiteID = data.WebsiteID
	case data.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:     comicID,
		model.DBComicGenericRID:         rid,
		model.DBWebsiteGenericWebsiteID: websiteID,
		model.DBComicCoverRelativeURL:   data.RelativeURL,
		model.DBComicCoverPriority:      data.Priority,
	})
	sql := "INSERT INTO " + model.DBComicCover + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicCoverRelativeURL
		sql += ", w." + model.DBComicCoverPriority
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicCoverSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicCoverSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicCover(ctx context.Context, conds any) (*model.ComicCover, error) {
	var result model.ComicCover
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicCoverRelativeURL
	sql += ", w." + model.DBComicCoverPriority
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBComicCover + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicCover(ctx context.Context, data model.SetComicCover, conds any, v *model.ComicCover) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.RID != nil {
		data0[model.DBComicGenericRID] = data.RID
	}
	switch {
	case data.WebsiteID != nil:
		data0[model.DBWebsiteGenericWebsiteID] = data.WebsiteID
	case data.WebsiteDomain != nil:
		data0[model.DBWebsiteGenericWebsiteID] = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	if data.RelativeURL != nil {
		data0[model.DBComicCoverRelativeURL] = data.RelativeURL
	}
	if data.Priority != nil {
		data0[model.DBComicCoverPriority] = data.Priority
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicCover + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicCoverRelativeURL
		sql += ", w." + model.DBComicCoverPriority
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicCoverSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicCoverSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicCover(ctx context.Context, conds any, v *model.ComicCover) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicCover + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicCoverRelativeURL
		sql += ", w." + model.DBComicCoverPriority
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
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

func (db Database) ListComicCover(ctx context.Context, params model.ListParams) ([]*model.ComicCover, error) {
	result := []*model.ComicCover{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicCoverRelativeURL
	sql += ", w." + model.DBComicCoverPriority
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBComicCover + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicGenericRID})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicCoverPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicCover(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicCover, conds)
}

func comicCoverSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicCoverFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicCoverFKey1:
				return model.GenericError("website does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists {
			switch errDatabase.Name {
			case NameErrComicCoverKey0:
				return model.GenericError("same comic id + rid already exists")
			case NameErrComicCoverKey1:
				return model.GenericError("same comic id + website id + relative url already exists")
			}
		}
	}
	return err
}

func (db Database) AddComicSynopsis(ctx context.Context, data model.AddComicSynopsis, v *model.ComicSynopsis) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var rid any
	switch {
	case data.RID != nil:
		rid = data.RID
	default:
		rid = utila.RandomString(utila.RandomStringGeneral, model.ComicGenericRIDLength)
	}
	var languageID any
	switch {
	case data.LanguageID != nil:
		languageID = data.LanguageID
	case data.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBComicGenericRID:           rid,
		model.DBLanguageGenericLanguageID: languageID,
		model.DBComicSynopsisSynopsis:     data.Synopsis,
		model.DBComicSynopsisVersion:      data.Version,
		model.DBComicSynopsisRomanized:    data.Romanized,
	})
	sql := "INSERT INTO " + model.DBComicSynopsis + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicSynopsisSynopsis
		sql += ", w." + model.DBComicSynopsisVersion + ", w." + model.DBComicSynopsisRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicSynopsisSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicSynopsisSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicSynopsis(ctx context.Context, conds any) (*model.ComicSynopsis, error) {
	var result model.ComicSynopsis
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicSynopsisSynopsis
	sql += ", w." + model.DBComicSynopsisVersion + ", w." + model.DBComicSynopsisRomanized
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComicSynopsis + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicSynopsis(ctx context.Context, data model.SetComicSynopsis, conds any, v *model.ComicSynopsis) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.RID != nil {
		data0[model.DBComicGenericRID] = data.RID
	}
	switch {
	case data.LanguageID != nil:
		data0[model.DBLanguageGenericLanguageID] = data.LanguageID
	case data.LanguageIETF != nil:
		data0[model.DBLanguageGenericLanguageID] = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	if data.Synopsis != nil {
		data0[model.DBComicSynopsisSynopsis] = data.Synopsis
	}
	if data.Version != nil {
		data0[model.DBComicSynopsisVersion] = data.Version
	}
	if data.Romanized != nil {
		data0[model.DBComicSynopsisRomanized] = data.Romanized
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicSynopsis + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicSynopsisSynopsis
		sql += ", w." + model.DBComicSynopsisVersion + ", w." + model.DBComicSynopsisRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicSynopsisSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicSynopsisSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicSynopsis(ctx context.Context, conds any, v *model.ComicSynopsis) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicSynopsis + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicSynopsisSynopsis
		sql += ", w." + model.DBComicSynopsisVersion + ", w." + model.DBComicSynopsisRomanized
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
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

func (db Database) ListComicSynopsis(ctx context.Context, params model.ListParams) ([]*model.ComicSynopsis, error) {
	result := []*model.ComicSynopsis{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBLanguageGenericLanguageID + ", w." + model.DBComicSynopsisSynopsis
	sql += ", w." + model.DBComicSynopsisVersion + ", w." + model.DBComicSynopsisRomanized
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBComicSynopsis + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicGenericRID})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicSynopsisPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicSynopsis(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicSynopsis, conds)
}

func comicSynopsisSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicSynopsisFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicSynopsisFKey1:
				return model.GenericError("language does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists {
			switch errDatabase.Name {
			case NameErrComicSynopsisKey0:
				return model.GenericError("same comic id + rid already exists")
			case NameErrComicSynopsisKey1:
				return model.GenericError("same comic id + synopsis already exists")
			}
		}
	}
	return err
}

func (db Database) AddComicExternal(ctx context.Context, data model.AddComicExternal, v *model.ComicExternal) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var rid any
	switch {
	case data.RID != nil:
		rid = data.RID
	default:
		rid = utila.RandomString(utila.RandomStringGeneral, model.ComicGenericRIDLength)
	}
	var websiteID any
	switch {
	case data.WebsiteID != nil:
		websiteID = data.WebsiteID
	case data.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:      comicID,
		model.DBComicGenericRID:          rid,
		model.DBWebsiteGenericWebsiteID:  websiteID,
		model.DBComicExternalRelativeURL: data.RelativeURL,
		model.DBComicExternalOfficial:    data.Official,
	})
	sql := "INSERT INTO " + model.DBComicExternal + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicExternalRelativeURL
		sql += ", w." + model.DBComicExternalOfficial
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicExternalSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicExternalSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicExternal(ctx context.Context, conds any) (*model.ComicExternal, error) {
	var result model.ComicExternal
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicExternalRelativeURL
	sql += ", w." + model.DBComicExternalOfficial
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBComicExternal + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicExternal(ctx context.Context, data model.SetComicExternal, conds any, v *model.ComicExternal) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.RID != nil {
		data0[model.DBComicGenericRID] = data.RID
	}
	switch {
	case data.WebsiteID != nil:
		data0[model.DBWebsiteGenericWebsiteID] = data.WebsiteID
	case data.WebsiteDomain != nil:
		data0[model.DBWebsiteGenericWebsiteID] = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	if data.RelativeURL != nil {
		data0[model.DBComicExternalRelativeURL] = data.RelativeURL
	}
	if data.Official != nil {
		data0[model.DBComicExternalOfficial] = data.Official
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicExternal + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicExternalRelativeURL
		sql += ", w." + model.DBComicExternalOfficial
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicExternalSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicExternalSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicExternal(ctx context.Context, conds any, v *model.ComicExternal) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicExternal + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicExternalRelativeURL
		sql += ", w." + model.DBComicExternalOfficial
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
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

func (db Database) ListComicExternal(ctx context.Context, params model.ListParams) ([]*model.ComicExternal, error) {
	result := []*model.ComicExternal{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicGenericRID
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBComicExternalRelativeURL
	sql += ", w." + model.DBComicExternalOfficial
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBComicExternal + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicGenericRID})
	}
	params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicExternalPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicExternal(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicExternal, conds)
}

func comicExternalSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicExternalFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicExternalFKey1:
				return model.GenericError("website does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists {
			switch errDatabase.Name {
			case NameErrComicExternalKey0:
				return model.GenericError("same comic id + rid already exists")
			case NameErrComicExternalKey1:
				return model.GenericError("same comic id + website id + relative url already exists")
			}
		}
	}
	return err
}

func (db Database) AddComicCategory(ctx context.Context, data model.AddComicCategory, v *model.ComicCategory) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var categoryID any
	switch {
	case data.CategoryID != nil:
		categoryID = data.CategoryID
	case data.CategoryCode != nil:
		categoryID = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.CategoryTypeID,
			TypeCode: data.CategoryTypeCode,
			Code:     *data.CategoryCode,
		})
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:       comicID,
		model.DBCategoryGenericCategoryID: categoryID,
	})
	sql := "INSERT INTO " + model.DBComicCategory + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBCategoryGenericCategoryID
		sql += ", l." + model.DBCategoryTypeID + " AS category_type_id"
		sql += ", l." + model.DBCategoryCode + " AS category_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryGenericCategoryID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicCategorySetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicCategorySetError(err)
		}
	}
	return nil
}

func (db Database) GetComicCategory(ctx context.Context, conds any) (*model.ComicCategory, error) {
	var result model.ComicCategory
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBCategoryGenericCategoryID
	sql += ", l." + model.DBCategoryTypeID + " AS category_type_id"
	sql += ", l." + model.DBCategoryCode + " AS category_code"
	sql += " FROM " + model.DBComicCategory + " w JOIN " + model.DBCategory + " l"
	sql += " ON w." + model.DBCategoryGenericCategoryID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicCategory(ctx context.Context, data model.SetComicCategory, conds any, v *model.ComicCategory) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	switch {
	case data.CategoryID != nil:
		data0[model.DBCategoryGenericCategoryID] = data.CategoryID
	case data.CategoryCode != nil:
		data0[model.DBCategoryGenericCategoryID] = model.DBCategorySIDToID(model.CategorySID{
			TypeID:   data.CategoryTypeID,
			TypeCode: data.CategoryTypeCode,
			Code:     *data.CategoryCode,
		})
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicCategory + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBCategoryGenericCategoryID
		sql += ", l." + model.DBCategoryTypeID + " AS category_type_id"
		sql += ", l." + model.DBCategoryCode + " AS category_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryGenericCategoryID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicCategorySetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicCategorySetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicCategory(ctx context.Context, conds any, v *model.ComicCategory) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicCategory + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBCategoryGenericCategoryID
		sql += ", l." + model.DBCategoryTypeID + " AS category_type_id"
		sql += ", l." + model.DBCategoryCode + " AS category_code"
		sql += " FROM data w JOIN " + model.DBCategory + " l"
		sql += " ON w." + model.DBCategoryGenericCategoryID + " = l." + model.DBGenericID
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

func (db Database) ListComicCategory(ctx context.Context, params model.ListParams) ([]*model.ComicCategory, error) {
	result := []*model.ComicCategory{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBCategoryGenericCategoryID
	sql += ", l." + model.DBCategoryTypeID + " AS category_type_id"
	sql += ", l." + model.DBCategoryCode + " AS category_code"
	sql += " FROM " + model.DBComicCategory + " w JOIN " + model.DBCategory + " l"
	sql += " ON w." + model.DBCategoryGenericCategoryID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBCategoryGenericCategoryID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicCategoryPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicCategory(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicCategory, conds)
}

func comicCategorySetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicCategoryFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicCategoryFKey1:
				return model.GenericError("category does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicCategoryPKey {
			return model.GenericError("same category id already exists")
		}
	}
	return err
}

func (db Database) AddComicTag(ctx context.Context, data model.AddComicTag, v *model.ComicTag) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var tagID any
	switch {
	case data.TagID != nil:
		tagID = data.TagID
	case data.TagCode != nil:
		tagID = model.DBTagSIDToID(model.TagSID{
			TypeID:   data.TagTypeID,
			TypeCode: data.TagTypeCode,
			Code:     *data.TagCode,
		})
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBTagGenericTagID:     tagID,
	})
	sql := "INSERT INTO " + model.DBComicTag + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBTagGenericTagID
		sql += ", l." + model.DBTagTypeID + " AS tag_type_id"
		sql += ", l." + model.DBTagCode + " AS tag_code"
		sql += " FROM data w JOIN " + model.DBTag + " l"
		sql += " ON w." + model.DBTagGenericTagID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicTagSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicTagSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicTag(ctx context.Context, conds any) (*model.ComicTag, error) {
	var result model.ComicTag
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBTagGenericTagID
	sql += ", l." + model.DBTagTypeID + " AS tag_type_id"
	sql += ", l." + model.DBTagCode + " AS tag_code"
	sql += " FROM " + model.DBComicTag + " w JOIN " + model.DBTag + " l"
	sql += " ON w." + model.DBTagGenericTagID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicTag(ctx context.Context, data model.SetComicTag, conds any, v *model.ComicTag) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	switch {
	case data.TagID != nil:
		data0[model.DBTagGenericTagID] = data.TagID
	case data.TagCode != nil:
		data0[model.DBTagGenericTagID] = model.DBTagSIDToID(model.TagSID{
			TypeID:   data.TagTypeID,
			TypeCode: data.TagTypeCode,
			Code:     *data.TagCode,
		})
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicTag + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBTagGenericTagID
		sql += ", l." + model.DBTagTypeID + " AS tag_type_id"
		sql += ", l." + model.DBTagCode + " AS tag_code"
		sql += " FROM data w JOIN " + model.DBTag + " l"
		sql += " ON w." + model.DBTagGenericTagID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicTagSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicTagSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicTag(ctx context.Context, conds any, v *model.ComicTag) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicTag + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBTagGenericTagID
		sql += ", l." + model.DBTagTypeID + " AS tag_type_id"
		sql += ", l." + model.DBTagCode + " AS tag_code"
		sql += " FROM data w JOIN " + model.DBTag + " l"
		sql += " ON w." + model.DBTagGenericTagID + " = l." + model.DBGenericID
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

func (db Database) ListComicTag(ctx context.Context, params model.ListParams) ([]*model.ComicTag, error) {
	result := []*model.ComicTag{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBTagGenericTagID
	sql += ", l." + model.DBTagTypeID + " AS tag_type_id"
	sql += ", l." + model.DBTagCode + " AS tag_code"
	sql += " FROM " + model.DBComicTag + " w JOIN " + model.DBTag + " l"
	sql += " ON w." + model.DBTagGenericTagID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBTagGenericTagID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicTagPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicTag(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicTag, conds)
}

func comicTagSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicTagFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicTagFKey1:
				return model.GenericError("tag does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicTagPKey {
			return model.GenericError("same tag id already exists")
		}
	}
	return err
}

func (db Database) AddComicRelationType(ctx context.Context, data model.AddComicRelationType, v *model.ComicRelationType) error {
	if err := db.GenericAdd(ctx, model.DBComicRelationType, map[string]any{
		model.DBComicRelationTypeCode: data.Code,
		model.DBComicRelationTypeName: data.Name,
	}, v); err != nil {
		return comicRelationTypeSetError(err)
	}
	return nil
}

func (db Database) GetComicRelationType(ctx context.Context, conds any) (*model.ComicRelationType, error) {
	var result model.ComicRelationType
	if err := db.GenericGet(ctx, model.DBComicRelationType, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicRelationType(ctx context.Context, data model.SetComicRelationType, conds any, v *model.ComicRelationType) error {
	data0 := map[string]any{}
	if data.Code != nil {
		data0[model.DBComicRelationTypeCode] = data.Code
	}
	if data.Name != nil {
		data0[model.DBComicRelationTypeName] = data.Name
	}
	if err := db.GenericUpdate(ctx, model.DBComicRelationType, data0, conds, v); err != nil {
		return comicRelationTypeSetError(err)
	}
	return nil
}

func (db Database) DeleteComicRelationType(ctx context.Context, conds any, v *model.ComicRelationType) error {
	return db.GenericDelete(ctx, model.DBComicRelationType, conds, v)
}

func (db Database) ListComicRelationType(ctx context.Context, params model.ListParams) ([]*model.ComicRelationType, error) {
	result := []*model.ComicRelationType{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicRelationTypeCode})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicRelationTypePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBComicRelationType, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicRelationType(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicRelationType, conds)
}

func comicRelationTypeSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicRelationTypeKey {
			return model.GenericError("same code already exists")
		}
	}
	return err
}

func (db Database) AddComicRelation(ctx context.Context, data model.AddComicRelation, v *model.ComicRelation) error {
	var typeID any
	switch {
	case data.TypeID != nil:
		typeID = data.TypeID
	case data.TypeCode != nil:
		typeID = model.DBComicRelationTypeCodeToID(*data.TypeCode)
	}
	var parentID any
	switch {
	case data.ParentID != nil:
		parentID = data.ParentID
	case data.ParentCode != nil:
		parentID = model.DBComicCodeToID(*data.ParentCode)
	}
	var childID any
	switch {
	case data.ChildID != nil:
		childID = data.ChildID
	case data.ChildCode != nil:
		childID = model.DBComicCodeToID(*data.ChildCode)
	}
	if v == nil {
		v = new(model.ComicRelation)
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
		model.DBComicRelationTypeID:   typeID,
		model.DBComicRelationParentID: parentID,
		model.DBComicRelationChildID:  childID,
	})
	sql := "INSERT INTO " + model.DBComicRelation + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicRelationParentID + ", w." + model.DBComicRelationTypeID
		sql += ", w." + model.DBComicRelationChildID + ", l." + model.DBComicCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicRelationChildID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
	}
	{
		var parentLoop *bool
		args := []any{}
		cte := "WITH RECURSIVE childs AS ("
		cte += " SELECT " + model.DBComicRelationChildID + ", " + model.DBComicRelationParentID
		cte += "  FROM " + model.DBComic + " WHERE " + model.DBComicRelationTypeID + " = " + utila.Utoa(v.TypeID)
		cte += " UNION SELECT childs." + model.DBComicRelationChildID + ", parents." + model.DBComicRelationParentID
		cte += "  FROM " + model.DBComic + " AS parents"
		cte += " JOIN childs ON childs." + model.DBComicRelationParentID + " = parents." + model.DBComicRelationChildID
		cte += ")"
		sql := cte + " SELECT (" + utila.Utoa(v.ParentID) + ", " + utila.Utoa(v.ChildID) + ") IN (SELECT * FROM childs)"
		if err := db.QueryOne(ctx, &parentLoop, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
		if parentLoop != nil && *parentLoop {
			return model.GenericError("comic relation loop detected")
		}
	}
	return db.ContextTransactionCommit(ctx)
}

func (db Database) GetComicRelation(ctx context.Context, conds any) (*model.ComicRelation, error) {
	var result model.ComicRelation
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicRelationParentID + ", w." + model.DBComicRelationTypeID
	sql += ", w." + model.DBComicRelationChildID + ", l." + model.DBComicCode + " AS child_code"
	sql += " FROM " + model.DBComicRelation + " w JOIN " + model.DBComic + " l"
	sql += " ON w." + model.DBComicRelationChildID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicRelation(ctx context.Context, data model.SetComicRelation, conds any, v *model.ComicRelation) error {
	data0 := map[string]any{}
	switch {
	case data.TypeID != nil:
		data0[model.DBComicRelationTypeID] = data.TypeID
	case data.TypeCode != nil:
		data0[model.DBComicRelationTypeID] = model.DBComicRelationTypeCodeToID(*data.TypeCode)
	}
	switch {
	case data.ParentID != nil:
		data0[model.DBComicRelationParentID] = data.ParentID
	case data.ParentCode != nil:
		data0[model.DBComicRelationParentID] = model.DBComicCodeToID(*data.ParentCode)
	}
	switch {
	case data.ChildID != nil:
		data0[model.DBComicRelationChildID] = data.ChildID
	case data.ChildCode != nil:
		data0[model.DBComicRelationChildID] = model.DBComicCodeToID(*data.ChildCode)
	}
	if v == nil {
		v = new(model.ComicRelation)
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
		sql += ", w." + model.DBComicRelationParentID + ", w." + model.DBComicRelationTypeID
		sql += ", w." + model.DBComicRelationChildID + ", l." + model.DBComicCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicRelationChildID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
	}
	{
		var parentLoop *bool
		args := []any{}
		cte := "WITH RECURSIVE childs AS ("
		cte += " SELECT " + model.DBComicRelationChildID + ", " + model.DBComicRelationParentID
		cte += "  FROM " + model.DBComic + " WHERE " + model.DBComicRelationTypeID + " = " + utila.Utoa(v.TypeID)
		cte += " UNION SELECT childs." + model.DBComicRelationChildID + ", parents." + model.DBComicRelationParentID
		cte += "  FROM " + model.DBComic + " AS parents"
		cte += " JOIN childs ON childs." + model.DBComicRelationParentID + " = parents." + model.DBComicRelationChildID
		cte += ")"
		sql := cte + " SELECT (" + utila.Utoa(v.ParentID) + ", " + utila.Utoa(v.ChildID) + ") IN (SELECT * FROM childs)"
		if err := db.QueryOne(ctx, &parentLoop, sql, args...); err != nil {
			return comicRelationSetError(err)
		}
		if parentLoop != nil && *parentLoop {
			return model.GenericError("comic relation loop detected")
		}
	}
	return db.ContextTransactionCommit(ctx)
}

func (db Database) DeleteComicRelation(ctx context.Context, conds any, v *model.ComicRelation) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicRelation + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicRelationParentID + ", w." + model.DBComicRelationTypeID
		sql += ", w." + model.DBComicRelationChildID + ", l." + model.DBComicCode + " AS child_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicRelationChildID + " = l." + model.DBGenericID
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

func (db Database) ListComicRelation(ctx context.Context, params model.ListParams) ([]*model.ComicRelation, error) {
	result := []*model.ComicRelation{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicRelationParentID + ", w." + model.DBComicRelationTypeID
	sql += ", w." + model.DBComicRelationChildID + ", l." + model.DBComicCode + " AS child_code"
	sql += " FROM " + model.DBComicRelation + " w JOIN " + model.DBComic + " l"
	sql += " ON w." + model.DBComicRelationChildID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicRelationChildID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicRelationPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicRelation(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicRelation, conds)
}

func comicRelationSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicRelationFKey0:
				return model.GenericError("comic relation type does not exist")
			case NameErrComicRelationFKey1:
				return model.GenericError("parent comic does not exist")
			case NameErrComicRelationFKey2:
				return model.GenericError("child comic does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicRelationPKey {
			return model.GenericError("same type id + child id already exists")
		}
		if errDatabase.Code == CodeErrValidation && errDatabase.Name == NameErrComicRelationCheck {
			return model.GenericError("parent comic and child comic cannot be same")
		}
	}
	return err
}

const (
	NameErrComicChapterFKey = "comic_chapter_comic_id_fkey"
	NameErrComicChapterKey  = "comic_chapter_comic_id_chapter_version_key"
)

func (db Database) AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	if err := db.GenericAdd(ctx, model.DBComicChapter, map[string]any{
		model.DBComicGenericComicID:    comicID,
		model.DBComicChapterChapter:    data.Chapter,
		model.DBComicChapterVersion:    data.Version,
		model.DBComicChapterVolume:     data.Volume,
		model.DBComicChapterReleasedAt: data.ReleasedAt,
	}, v); err != nil {
		return comicChapterSetError(err)
	}
	return nil
}

func (db Database) GetComicChapter(ctx context.Context, conds any) (*model.ComicChapter, error) {
	var result model.ComicChapter
	if err := db.GenericGet(ctx, model.DBComicChapter, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.Chapter != nil {
		data0[model.DBComicChapterChapter] = data.Chapter
	}
	if data.Version != nil {
		data0[model.DBComicChapterVersion] = data.Version
	}
	if data.Volume != nil {
		data0[model.DBComicChapterVolume] = data.Volume
	}
	if data.ReleasedAt != nil {
		data0[model.DBComicChapterReleasedAt] = data.ReleasedAt
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	if err := db.GenericUpdate(ctx, model.DBComicChapter, data0, conds, v); err != nil {
		return comicChapterSetError(err)
	}
	return nil
}

func (db Database) DeleteComicChapter(ctx context.Context, conds any, v *model.ComicChapter) error {
	return db.GenericDelete(ctx, model.DBComicChapter, conds, v)
}

func (db Database) ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error) {
	result := []*model.ComicChapter{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicChapterReleasedAt})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicChapterPaginationDef}
	}
	if err := db.GenericList(ctx, model.DBComicChapter, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicChapter(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicChapter, conds)
}

func comicChapterSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign && errDatabase.Name == NameErrComicChapterFKey {
			return model.GenericError("comic does not exist")
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicChapterKey {
			return model.GenericError("same comic id + chapter + version already exists")
		}
	}
	return err
}
