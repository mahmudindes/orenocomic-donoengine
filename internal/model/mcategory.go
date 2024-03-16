package model

import (
	"strconv"
	"time"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
)

func init() {
	CategoryTypeOrderByAllow = append(CategoryTypeOrderByAllow, GenericOrderByAllow...)
}

const (
	CategoryTypeCodeMax       = 24
	CategoryTypeNameMax       = 24
	CategoryTypeOrderBysMax   = 3
	CategoryTypePaginationDef = 10
	CategoryTypePaginationMax = 50
	DBCategoryType            = donoengine.ID + "." + "category_type"
	DBCategoryTypeCode        = "code"
	DBCategoryTypeName        = "name"
)

var (
	CategoryTypeOrderByAllow = []string{
		DBCategoryTypeCode,
		DBCategoryTypeName,
	}

	DBCategoryTypeCodeToID = func(code string) DBQueryValue {
		return DBQueryValue{
			Table:      DBCategoryType,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBCategoryTypeCode, Value: code},
		}
	}
)

type (
	CategoryType struct {
		ID        uint       `json:"id"`
		Code      string     `json:"code"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	AddCategoryType struct {
		Code string
		Name string
	}

	SetCategoryType struct {
		Code *string
		Name *string
	}
)

func (m AddCategoryType) Validate() error {
	return (SetCategoryType{
		Code: &m.Code,
		Name: &m.Name,
	}).Validate()
}

func (m SetCategoryType) Validate() error {
	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) > CategoryTypeCodeMax {
			max := strconv.FormatInt(CategoryTypeCodeMax, 10)
			return GenericError("code must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > CategoryTypeNameMax {
			max := strconv.FormatInt(CategoryTypeNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

func init() {
	CategoryOrderByAllow = append(CategoryOrderByAllow, GenericOrderByAllow...)
}

const (
	CategoryCodeMax       = 32
	CategoryNameMax       = 32
	CategoryOrderBysMax   = 3
	CategoryPaginationDef = 10
	CategoryPaginationMax = 50
	DBCategory            = donoengine.ID + "." + "category"
	DBCategoryTypeID      = "type_id"
	DBCategoryCode        = "code"
	DBCategoryName        = "name"
)

var (
	CategoryOrderByAllow = []string{
		DBCategoryCode,
		DBCategoryName,
	}

	DBCategorySIDToID = func(sid CategorySID) DBQueryValue {
		var typeID any
		switch {
		case sid.TypeID != nil:
			typeID = sid.TypeID
		case sid.TypeCode != nil:
			typeID = DBCategoryTypeCodeToID(*sid.TypeCode)
		}
		return DBQueryValue{
			Table:      DBCategory,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: map[string]any{
				DBCategoryTypeID: typeID,
				DBCategoryCode:   sid.Code,
			},
		}
	}
)

type (
	Category struct {
		ID        uint                `json:"id"`
		TypeID    uint                `json:"typeID"`
		TypeCode  string              `json:"typeCode"`
		Code      string              `json:"code"`
		Name      string              `json:"name"`
		Relations []*CategoryRelation `db:"-" json:"relations"`
		CreatedAt time.Time           `json:"createdAt"`
		UpdatedAt *time.Time          `json:"updatedAt"`
	}

	AddCategory struct {
		TypeID   *uint
		TypeCode *string
		Code     string
		Name     string
	}

	SetCategory struct {
		TypeID   *uint
		TypeCode *string
		Code     *string
		Name     *string
	}

	CategorySID struct {
		TypeID   *uint
		TypeCode *string
		Code     string
	}
)

func (m AddCategory) Validate() error {
	if m.TypeID == nil && m.TypeCode == nil {
		return GenericError("either category type id or category type code must exist")
	}

	return (SetCategory{
		TypeID:   m.TypeID,
		TypeCode: m.TypeCode,
		Code:     &m.Code,
		Name:     &m.Name,
	}).Validate()
}

func (m SetCategory) Validate() error {
	if err := (SetCategoryType{Code: m.TypeCode}).Validate(); err != nil {
		return GenericError("type " + err.Error())
	}

	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) > CategoryCodeMax {
			max := strconv.FormatInt(CategoryCodeMax, 10)
			return GenericError("code must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > CategoryNameMax {
			max := strconv.FormatInt(CategoryNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

func init() {
	CategoryRelationOrderByAllow = append(CategoryRelationOrderByAllow, GenericOrderByAllow...)
}

const (
	DBCategoryGenericCategoryID   = "category_id"
	CategoryRelationOrderBysMax   = 3
	CategoryRelationPaginationDef = 10
	CategoryRelationPaginationMax = 50
	DBCategoryRelation            = donoengine.ID + "." + "category_relation"
	DBCategoryRelationParentID    = "parent_id"
	DBCategoryRelationChildID     = "child_id"
)

var CategoryRelationOrderByAllow = []string{
	DBCategoryRelationChildID,
}

type (
	CategoryRelation struct {
		ParentID  uint       `json:"-"`
		ChildID   uint       `json:"categoryID"`
		ChildCode string     `json:"categoryCode"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	AddCategoryRelation struct {
		TypeID     *uint
		TypeCode   *string
		ParentID   *uint
		ParentCode *string
		ChildID    *uint
		ChildCode  *string
	}
	SetCategoryRelation struct {
		TypeID     *uint
		TypeCode   *string
		ParentID   *uint
		ParentCode *string
		ChildID    *uint
		ChildCode  *string
	}
	CategoryRelationSID struct {
		TypeID     *uint
		TypeCode   *string
		ParentID   *uint
		ParentCode *string
		ChildID    *uint
		ChildCode  *string
	}
)

func (m AddCategoryRelation) Validate() error {
	if m.ParentID == nil && m.ParentCode == nil {
		return GenericError("either parent category id or parent category code must exist")
	}

	if m.ChildID == nil && m.ChildCode == nil {
		return GenericError("either child category id or child category code must exist")
	}

	return (&SetComicRelation{
		TypeID:     m.TypeID,
		TypeCode:   m.TypeCode,
		ParentID:   m.ParentID,
		ParentCode: m.ParentCode,
		ChildID:    m.ChildID,
		ChildCode:  m.ChildCode,
	}).Validate()
}
func (m SetCategoryRelation) Validate() error {
	if m.TypeID == nil && m.TypeCode == nil {
		return GenericError("either category relation type id or category relation type code must exist")
	}

	if err := (SetCategory{Code: m.ParentCode}).Validate(); err != nil {
		return GenericError("parent category " + err.Error())
	}

	if err := (SetCategory{Code: m.ChildCode}).Validate(); err != nil {
		return GenericError("child category " + err.Error())
	}

	return nil
}
