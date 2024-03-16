package model

import (
	"strconv"
	"time"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

func init() {
	LanguageOrderByAllow = append(LanguageOrderByAllow, GenericOrderByAllow...)
}

const (
	LanguageIETFMax       = 12
	LanguageNameMax       = 24
	LanguageOrderBysMax   = 3
	LanguagePaginationDef = 10
	LanguagePaginationMax = 50
	DBLanguage            = donoengine.ID + "." + "language"
	DBLanguageIETF        = "ietf"
	DBLanguageName        = "name"
)

var (
	LanguageOrderByAllow = []string{
		DBLanguageIETF,
		DBLanguageName,
	}

	DBLanguageIETFToID = func(ietf string) DBQueryValue {
		return DBQueryValue{
			Table:      DBLanguage,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBLanguageIETF, Value: ietf},
		}
	}
)

type (
	Language struct {
		ID        uint       `json:"id"`
		IETF      string     `json:"ietf"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	AddLanguage struct {
		IETF string
		Name string
	}

	SetLanguage struct {
		IETF *string `validate:"omitnil,bcp47_language_tag"`
		Name *string
	}
)

func (m AddLanguage) Validate() error {
	return (SetLanguage{
		IETF: &m.IETF,
		Name: &m.Name,
	}).Validate()
}

func (m SetLanguage) Validate() error {
	if m.IETF != nil {
		if *m.IETF == "" {
			return GenericError("ietf cannot be empty")
		}

		if len(*m.IETF) > LanguageIETFMax {
			max := strconv.FormatInt(LanguageIETFMax, 10)
			return GenericError("ietf must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > LanguageNameMax {
			max := strconv.FormatInt(LanguageNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	if err := utila.Validator.Struct(&m); err != nil {
		return GenericError("language data is not valid")
	}

	return nil
}

const DBLanguageGenericLanguageID = "language_id"
