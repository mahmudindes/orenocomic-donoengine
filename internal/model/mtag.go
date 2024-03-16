package model

import (
	"strconv"
	"time"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
)

func init() {
	TagTypeOrderByAllow = append(TagTypeOrderByAllow, GenericOrderByAllow...)
}

const (
	TagTypeCodeMax       = 24
	TagTypeNameMax       = 24
	TagTypeOrderBysMax   = 3
	TagTypePaginationDef = 10
	TagTypePaginationMax = 50
	DBTagType            = donoengine.ID + "." + "tag_type"
	DBTagTypeCode        = "code"
	DBTagTypeName        = "name"
)

var (
	TagTypeOrderByAllow = []string{
		DBCategoryTypeCode,
		DBCategoryTypeName,
	}

	DBTagTypeCodeToID = func(code string) DBQueryValue {
		return DBQueryValue{
			Table:      DBTagType,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBTagTypeCode, Value: code},
		}
	}
)

type (
	TagType struct {
		ID        uint       `json:"id"`
		Code      string     `json:"code"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	AddTagType struct {
		Code string
		Name string
	}

	SetTagType struct {
		Code *string
		Name *string
	}
)

func (m AddTagType) Validate() error {
	return (SetCategoryType{
		Code: &m.Code,
		Name: &m.Name,
	}).Validate()
}

func (m SetTagType) Validate() error {
	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) > TagTypeCodeMax {
			max := strconv.FormatInt(TagTypeCodeMax, 10)
			return GenericError("code must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > TagTypeNameMax {
			max := strconv.FormatInt(TagTypeNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

func init() {
	TagOrderByAllow = append(TagOrderByAllow, GenericOrderByAllow...)
}

const (
	TagCodeMax       = 32
	TagNameMax       = 32
	TagOrderBysMax   = 3
	TagPaginationDef = 10
	TagPaginationMax = 50
	DBTag            = donoengine.ID + "." + "tag"
	DBTagTypeID      = "type_id"
	DBTagCode        = "code"
	DBTagName        = "name"
)

var (
	TagOrderByAllow = []string{
		DBTagCode,
		DBTagName,
	}

	DBTagSIDToID = func(sid TagSID) DBQueryValue {
		var typeID any
		switch {
		case sid.TypeID != nil:
			typeID = sid.TypeID
		case sid.TypeCode != nil:
			typeID = DBTagTypeCodeToID(*sid.TypeCode)
		}
		return DBQueryValue{
			Table:      DBTag,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: map[string]any{
				DBTagTypeID: typeID,
				DBTagCode:   sid.Code,
			},
		}
	}
)

type (
	Tag struct {
		ID        uint       `json:"id"`
		TypeID    uint       `json:"typeID"`
		TypeCode  string     `json:"typeCode"`
		Code      string     `json:"code"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	AddTag struct {
		TypeID   *uint
		TypeCode *string
		Code     string
		Name     string
	}

	SetTag struct {
		TypeID   *uint
		TypeCode *string
		Code     *string
		Name     *string
	}

	TagSID struct {
		TypeID   *uint
		TypeCode *string
		Code     string
	}
)

func (m AddTag) Validate() error {
	if m.TypeID == nil && m.TypeCode == nil {
		return GenericError("either tag type id or tag type code must exist")
	}

	return (SetTag{
		TypeID:   m.TypeID,
		TypeCode: m.TypeCode,
		Code:     &m.Code,
		Name:     &m.Name,
	}).Validate()
}

func (m SetTag) Validate() error {
	if err := (SetTagType{Code: m.TypeCode}).Validate(); err != nil {
		return GenericError("type " + err.Error())
	}

	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) > TagCodeMax {
			max := strconv.FormatInt(TagCodeMax, 10)
			return GenericError("code must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > TagNameMax {
			max := strconv.FormatInt(TagNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

const DBTagGenericTagID = "tag_id"
