package model

import (
	"slices"
	"strconv"
	"time"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
)

func init() {
	ComicOrderByAllow = append(ComicOrderByAllow, GenericOrderByAllow...)
}

const (
	ComicCodeLength      = 8
	ComicOrderBysMax     = 5
	ComicPaginationDef   = 10
	ComicPaginationMax   = 30
	DBComic              = donoengine.ID + "." + "comic"
	DBComicCode          = "code"
	DBComicPublishedFrom = "published_from"
	DBComicPublishedTo   = "published_to"
	DBComicTotalChapter  = "total_chapter"
	DBComicTotalVolume   = "total_volume"
	DBComicNSFW          = "nsfw"
	DBComicNSFL          = "nsfl"
	DBComicAdditionals   = "additionals"
)

var (
	ComicOrderByAllow = []string{
		DBComicCode,
		DBComicPublishedFrom,
		DBComicPublishedTo,
		DBComicTotalChapter,
		DBComicTotalVolume,
		DBComicNSFW,
		DBComicNSFL,
		DBLanguageGenericLanguageID,
	}

	ComicSetNullAllow = []string{
		DBComicPublishedFrom,
		DBComicPublishedTo,
		DBComicTotalChapter,
		DBComicTotalVolume,
		DBComicNSFW,
		DBComicNSFL,
		DBLanguageGenericLanguageID,
		DBComicAdditionals,
	}

	DBComicCodeToID = func(code string) DBQueryValue {
		return DBQueryValue{
			Table:      DBComic,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBComicCode, Value: code},
		}
	}
)

type (
	Comic struct {
		ID            uint             `json:"id"`
		Code          string           `json:"code"`
		Titles        []*ComicTitle    `db:"-" json:"titles"`
		Covers        []*ComicCover    `db:"-" json:"covers"`
		Synopses      []*ComicSynopsis `db:"-" json:"synopses"`
		PublishedFrom *time.Time       `json:"publishedFrom"`
		PublishedTo   *time.Time       `json:"publishedTo"`
		TotalChapter  *int             `json:"totalChapter"`
		TotalVolume   *int             `json:"totalVolume"`
		NSFW          *int             `json:"nsfw"`
		NSFL          *int             `json:"nsfl"`
		LanguageID    *uint            `json:"languageID"`
		LanguageIETF  *string          `json:"languageIETF"`
		Chapters      []*ComicChapter  `db:"-" json:"chapters"`
		Externals     []*ComicExternal `db:"-" json:"externals"`
		Categories    []*Category      `db:"-" json:"categories"`
		Tags          []*Tag           `db:"-" json:"tags"`
		Relations     []*ComicRelation `db:"-" json:"relations"`
		Additionals   map[string]any   `json:"additionals"`
		CreatedAt     time.Time        `json:"createdAt"`
		UpdatedAt     *time.Time       `json:"updatedAt"`
	}

	AddComic struct {
		Code          *string
		LanguageID    *uint
		LanguageIETF  *string
		PublishedFrom *time.Time
		PublishedTo   *time.Time
		TotalChapter  *int
		TotalVolume   *int
		NSFW          *int
		NSFL          *int
		Additionals   map[string]any
	}

	SetComic struct {
		Code          *string
		LanguageID    *uint
		LanguageIETF  *string
		PublishedFrom *time.Time
		PublishedTo   *time.Time
		TotalChapter  *int
		TotalVolume   *int
		NSFW          *int
		NSFL          *int
		Additionals   map[string]any
		SetNull       []string
	}
)

func (m AddComic) Validate() error {
	return (SetComic{
		Code:          m.Code,
		LanguageID:    m.LanguageID,
		LanguageIETF:  m.LanguageIETF,
		PublishedFrom: m.PublishedFrom,
		PublishedTo:   m.PublishedTo,
		TotalChapter:  m.TotalChapter,
		TotalVolume:   m.TotalVolume,
		NSFW:          m.NSFW,
		NSFL:          m.NSFL,
		Additionals:   m.Additionals,
	}).Validate()
}

func (m SetComic) Validate() error {
	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) != ComicCodeLength {
			length := strconv.Itoa(ComicCodeLength)
			return GenericError("code must be " + length + " characters long")
		}
	}

	if m.LanguageIETF != nil {
		if err := (SetLanguage{IETF: m.LanguageIETF}).Validate(); err != nil {
			return GenericError("language " + err.Error())
		}
	}

	if m.PublishedFrom != nil && m.PublishedTo != nil {
		if m.PublishedFrom.After(*m.PublishedTo) {
			return GenericError("published from is after published to")
		}
	}

	if m.NSFW != nil && *m.NSFW < -1 && *m.NSFW > 1 {
		return GenericError("nsfw must be at last -1 and at most 1")
	}

	if m.NSFL != nil && *m.NSFL < -1 && *m.NSFL > 1 {
		return GenericError("nsfl must be at last -1 and at most 1")
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func init() {
	ComicTitleOrderByAllow = append(ComicTitleOrderByAllow, GenericOrderByAllow...)
	ComicCoverOrderByAllow = append(ComicCoverOrderByAllow, GenericOrderByAllow...)
	ComicSynopsisOrderByAllow = append(ComicSynopsisOrderByAllow, GenericOrderByAllow...)
	ComicExternalOrderByAllow = append(ComicExternalOrderByAllow, GenericOrderByAllow...)
	ComicCategoryOrderByAllow = append(ComicCategoryOrderByAllow, GenericOrderByAllow...)
	ComicTagOrderByAllow = append(ComicTagOrderByAllow, GenericOrderByAllow...)
	ComicRelationTypeOrderByAllow = append(ComicRelationTypeOrderByAllow, GenericOrderByAllow...)
	ComicRelationOrderByAllow = append(ComicRelationOrderByAllow, GenericOrderByAllow...)
}

const (
	ComicGenericRIDLength          = 4
	DBComicGenericComicID          = "comic_id"
	DBComicGenericRID              = "rid"
	ComicTitleTitleMax             = 255
	ComicTitleOrderBysMax          = 3
	ComicTitlePaginationDef        = 10
	ComicTitlePaginationMax        = 50
	DBComicTitle                   = donoengine.ID + "." + "comic_title"
	DBComicTitleTitle              = "title"
	DBComicTitleSynonym            = "synonym"
	DBComicTitleRomanized          = "romanized"
	ComicCoverRelativeURLMax       = 128
	ComicCoverOrderBysMax          = 3
	ComicCoverPaginationDef        = 10
	ComicCoverPaginationMax        = 50
	DBComicCover                   = donoengine.ID + "." + "comic_cover"
	DBComicCoverRelativeURL        = "relative_url"
	DBComicCoverPriority           = "priority"
	ComicSynopsisSynopsisMax       = 2048
	ComicSynopsisVersionMax        = 12
	ComicSynopsisOrderBysMax       = 3
	ComicSynopsisPaginationDef     = 10
	ComicSynopsisPaginationMax     = 50
	DBComicSynopsis                = donoengine.ID + "." + "comic_synopsis"
	DBComicSynopsisSynopsis        = "synopsis"
	DBComicSynopsisVersion         = "version"
	DBComicSynopsisRomanized       = "romanized"
	ComicExternalRelativeURLMax    = 128
	ComicExternalOrderBysMax       = 3
	ComicExternalPaginationDef     = 10
	ComicExternalPaginationMax     = 50
	DBComicExternal                = donoengine.ID + "." + "comic_external"
	DBComicExternalRelativeURL     = "relative_url"
	DBComicExternalOfficial        = "official"
	ComicCategoryOrderBysMax       = 3
	ComicCategoryPaginationDef     = 10
	ComicCategoryPaginationMax     = 50
	DBComicCategory                = donoengine.ID + "." + "comic_category"
	ComicTagOrderBysMax            = 3
	ComicTagPaginationDef          = 10
	ComicTagPaginationMax          = 50
	DBComicTag                     = donoengine.ID + "." + "comic_tag"
	ComicRelationTypeCodeMax       = 24
	ComicRelationTypeNameMax       = 24
	ComicRelationTypeOrderBysMax   = 3
	ComicRelationTypePaginationDef = 10
	ComicRelationTypePaginationMax = 50
	DBComicRelationType            = donoengine.ID + "." + "comic_relation_type"
	DBComicRelationTypeCode        = "code"
	DBComicRelationTypeName        = "name"
	ComicRelationOrderBysMax       = 3
	ComicRelationPaginationDef     = 10
	ComicRelationPaginationMax     = 50
	DBComicRelation                = donoengine.ID + "." + "comic_relation"
	DBComicRelationTypeID          = "type_id"
	DBComicRelationParentID        = "parent_id"
	DBComicRelationChildID         = "child_id"
)

var (
	ComicTitleOrderByAllow = []string{
		DBComicGenericRID,
		DBLanguageGenericLanguageID,
		DBComicTitleTitle,
		DBComicTitleSynonym,
		DBComicTitleRomanized,
	}
	ComicTitleSetNullAllow = []string{
		DBComicTitleSynonym,
		DBComicTitleRomanized,
	}

	ComicCoverOrderByAllow = []string{
		DBComicGenericRID,
		DBWebsiteGenericWebsiteID,
		DBComicCoverRelativeURL,
		DBComicCoverPriority,
	}
	ComicCoverSetNullAllow = []string{
		DBComicCoverPriority,
	}

	ComicSynopsisOrderByAllow = []string{
		DBComicGenericRID,
		DBLanguageGenericLanguageID,
		DBComicSynopsisSynopsis,
		DBComicSynopsisVersion,
		DBComicSynopsisRomanized,
	}
	ComicSynopsisSetNullAllow = []string{
		DBComicSynopsisVersion,
		DBComicSynopsisRomanized,
	}

	ComicExternalOrderByAllow = []string{
		DBComicGenericRID,
		DBWebsiteGenericWebsiteID,
		DBWebsiteGenericWebsiteDomain,
		DBComicExternalRelativeURL,
		DBComicExternalOfficial,
	}
	ComicExternalSetNullAllow = []string{
		DBComicExternalRelativeURL,
		DBComicExternalOfficial,
	}

	ComicCategoryOrderByAllow = []string{
		DBCategoryGenericCategoryID,
	}

	ComicTagOrderByAllow = []string{
		DBTagGenericTagID,
	}

	ComicRelationTypeOrderByAllow = []string{
		DBComicRelationTypeCode,
		DBComicRelationTypeName,
	}
	DBComicRelationTypeCodeToID = func(code string) DBQueryValue {
		return DBQueryValue{
			Table:      DBComicRelationType,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBComicRelationTypeCode, Value: code},
		}
	}

	ComicRelationOrderByAllow = []string{
		DBComicRelationTypeID,
		DBComicRelationChildID,
	}
)

type (
	ComicGenericSID struct {
		ComicID   *uint
		ComicCode *string
		RID       string
	}

	ComicTitle struct {
		ID           uint       `json:"id"`
		ComicID      uint       `json:"-"`
		RID          string     `json:"rid"`
		LanguageID   uint       `json:"languageID"`
		LanguageIETF string     `json:"languageIETF"`
		Title        string     `json:"title"`
		Synonym      *bool      `json:"synonym"`
		Romanized    *bool      `json:"romanized"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    *time.Time `json:"updatedAt"`
	}
	AddComicTitle struct {
		ComicID      *uint
		ComicCode    *string
		RID          *string
		LanguageID   *uint
		LanguageIETF *string
		Title        string
		Synonym      *bool
		Romanized    *bool
	}
	SetComicTitle struct {
		ComicID      *uint
		ComicCode    *string
		RID          *string
		LanguageID   *uint
		LanguageIETF *string
		Title        *string
		Synonym      *bool
		Romanized    *bool
		SetNull      []string
	}

	ComicCover struct {
		ID            uint       `json:"id"`
		ComicID       uint       `json:"-"`
		RID           string     `json:"rid"`
		WebsiteID     uint       `json:"websiteID"`
		WebsiteDomain string     `json:"websiteDomain"`
		RelativeURL   string     `json:"relativeURL"`
		Priority      *int       `json:"priority"`
		CreatedAt     time.Time  `json:"createdAt"`
		UpdatedAt     *time.Time `json:"updatedAt"`
	}
	AddComicCover struct {
		ComicID       *uint
		ComicCode     *string
		RID           *string
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   string
		Priority      *int
	}
	SetComicCover struct {
		ComicID       *uint
		ComicCode     *string
		RID           *string
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   *string
		Priority      *int
		SetNull       []string
	}

	ComicSynopsis struct {
		ID           uint       `json:"id"`
		ComicID      uint       `json:"-"`
		RID          string     `json:"rid"`
		LanguageID   uint       `json:"languageID"`
		LanguageIETF string     `json:"languageIETF"`
		Synopsis     string     `json:"synopsis"`
		Version      *string    `json:"version"`
		Romanized    *bool      `json:"romanized"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    *time.Time `json:"updatedAt"`
	}
	AddComicSynopsis struct {
		ComicID      *uint
		ComicCode    *string
		RID          *string
		LanguageID   *uint
		LanguageIETF *string
		Synopsis     string
		Version      *string
		Romanized    *bool
	}
	SetComicSynopsis struct {
		ComicID      *uint
		ComicCode    *string
		RID          *string
		LanguageID   *uint
		LanguageIETF *string
		Synopsis     *string
		Version      *string
		Romanized    *bool
		SetNull      []string
	}

	ComicExternal struct {
		ID            uint       `json:"id"`
		ComicID       uint       `json:"-"`
		RID           string     `json:"rid"`
		WebsiteID     uint       `json:"websiteID"`
		WebsiteDomain string     `json:"websiteDomain"`
		RelativeURL   *string    `json:"relativeURL"`
		Official      *bool      `json:"official"`
		CreatedAt     time.Time  `json:"createdAt"`
		UpdatedAt     *time.Time `json:"updatedAt"`
	}
	AddComicExternal struct {
		ComicID       *uint
		ComicCode     *string
		RID           *string
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   *string
		Official      *bool
	}
	SetComicExternal struct {
		ComicID       *uint
		ComicCode     *string
		RID           *string
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   *string
		Official      *bool
		SetNull       []string
	}

	ComicCategory struct {
		ComicID        uint       `json:"-"`
		CategoryID     uint       `json:"categoryID"`
		CategoryTypeID uint       `json:"categoryTypeID"`
		CategoryCode   string     `json:"categoryCode"`
		CreatedAt      time.Time  `json:"createdAt"`
		UpdatedAt      *time.Time `json:"updatedAt"`
	}
	AddComicCategory struct {
		ComicID          *uint
		ComicCode        *string
		CategoryID       *uint
		CategoryTypeID   *uint
		CategoryTypeCode *string
		CategoryCode     *string
	}
	SetComicCategory struct {
		ComicID          *uint
		ComicCode        *string
		CategoryID       *uint
		CategoryTypeID   *uint
		CategoryTypeCode *string
		CategoryCode     *string
	}
	ComicCategorySID struct {
		ComicID     *uint
		ComicCode   *string
		CategoryID  *uint
		CategorySID *CategorySID
	}

	ComicTag struct {
		ComicID   uint       `json:"-"`
		TagID     uint       `json:"tagID"`
		TagTypeID uint       `json:"tagTypeID"`
		TagCode   string     `json:"tagCode"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	AddComicTag struct {
		ComicID     *uint
		ComicCode   *string
		TagID       *uint
		TagTypeID   *uint
		TagTypeCode *string
		TagCode     *string
	}
	SetComicTag struct {
		ComicID     *uint
		ComicCode   *string
		TagID       *uint
		TagTypeID   *uint
		TagTypeCode *string
		TagCode     *string
	}
	ComicTagSID struct {
		ComicID   *uint
		ComicCode *string
		TagID     *uint
		TagSID    *TagSID
	}

	ComicRelationType struct {
		ID        uint       `json:"id"`
		Code      string     `json:"code"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	AddComicRelationType struct {
		Code string
		Name string
	}
	SetComicRelationType struct {
		Code *string
		Name *string
	}

	ComicRelation struct {
		ParentID  uint       `json:"-"`
		TypeID    uint       `json:"typeID"`
		ChildID   uint       `json:"comicID"`
		ChildCode string     `json:"comicCode"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	AddComicRelation struct {
		ParentID   *uint
		ParentCode *string
		TypeID     *uint
		TypeCode   *string
		ChildID    *uint
		ChildCode  *string
	}
	SetComicRelation struct {
		ParentID   *uint
		ParentCode *string
		TypeID     *uint
		TypeCode   *string
		ChildID    *uint
		ChildCode  *string
	}
	ComicRelationSID struct {
		ParentID   *uint
		ParentCode *string
		TypeID     *uint
		TypeCode   *string
		ChildID    *uint
		ChildCode  *string
	}
)

func (m AddComicTitle) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.LanguageID == nil && m.LanguageIETF == nil {
		return GenericError("either language id or language ietf must exist")
	}

	return (SetComicTitle{
		ComicID:      m.ComicID,
		ComicCode:    m.ComicCode,
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
		Title:        &m.Title,
		Synonym:      m.Synonym,
		Romanized:    m.Romanized,
	}).Validate()
}
func (m SetComicTitle) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.RID != nil {
		if *m.RID == "" {
			return GenericError("rid cannot be empty")
		}

		if len(*m.RID) != ComicGenericRIDLength {
			length := strconv.Itoa(ComicGenericRIDLength)
			return GenericError("rid must be " + length + " characters long")
		}
	}

	if err := (SetLanguage{IETF: m.LanguageIETF}).Validate(); err != nil {
		return GenericError("language " + err.Error())
	}

	if m.Title != nil {
		if *m.Title == "" {
			return GenericError("title cannot be empty")
		}

		if len(*m.Title) > ComicTitleTitleMax {
			max := strconv.FormatInt(ComicTitleTitleMax, 10)
			return GenericError("title must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func (m AddComicCover) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.WebsiteID == nil && m.WebsiteDomain == nil {
		return GenericError("either website id or website domain must exist")
	}

	return (SetComicCover{
		ComicID:       m.ComicID,
		ComicCode:     m.ComicCode,
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   &m.RelativeURL,
		Priority:      m.Priority,
	}).Validate()
}
func (m SetComicCover) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.RID != nil {
		if *m.RID == "" {
			return GenericError("rid cannot be empty")
		}

		if len(*m.RID) != ComicGenericRIDLength {
			length := strconv.Itoa(ComicGenericRIDLength)
			return GenericError("rid must be " + length + " characters long")
		}
	}

	if err := (SetWebsite{Domain: m.WebsiteDomain}).Validate(); err != nil {
		return GenericError("website " + err.Error())
	}

	if m.RelativeURL != nil {
		if *m.RelativeURL == "" {
			return GenericError("relative url cannot be empty")
		}

		if len(*m.RelativeURL) > ComicCoverRelativeURLMax {
			max := strconv.FormatInt(ComicCoverRelativeURLMax, 10)
			return GenericError("relative url must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicCoverSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func (m AddComicSynopsis) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.LanguageID == nil && m.LanguageIETF == nil {
		return GenericError("either language id or language ietf must exist")
	}

	return (SetComicSynopsis{
		ComicID:      m.ComicID,
		ComicCode:    m.ComicCode,
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
		Synopsis:     &m.Synopsis,
		Version:      m.Version,
		Romanized:    m.Romanized,
	}).Validate()
}
func (m SetComicSynopsis) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.RID != nil {
		if *m.RID == "" {
			return GenericError("rid cannot be empty")
		}

		if len(*m.RID) != ComicGenericRIDLength {
			length := strconv.Itoa(ComicGenericRIDLength)
			return GenericError("rid must be " + length + " characters long")
		}
	}

	if err := (SetLanguage{IETF: m.LanguageIETF}).Validate(); err != nil {
		return GenericError("language " + err.Error())
	}

	if m.Synopsis != nil {
		if *m.Synopsis == "" {
			return GenericError("synopsis cannot be empty")
		}

		if len(*m.Synopsis) > ComicSynopsisSynopsisMax {
			max := strconv.FormatInt(ComicSynopsisSynopsisMax, 10)
			return GenericError("synopsis must be at most " + max + " characters long")
		}
	}

	if m.Version != nil {
		if *m.Version == "" {
			return GenericError("version cannot be empty")
		}

		if len(*m.Version) > ComicSynopsisVersionMax {
			max := strconv.FormatInt(ComicSynopsisVersionMax, 10)
			return GenericError("version must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicSynopsisSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func (m AddComicExternal) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.WebsiteID == nil && m.WebsiteDomain == nil {
		return GenericError("either website id or website domain must exist")
	}

	return (SetComicExternal{
		ComicID:       m.ComicID,
		ComicCode:     m.ComicCode,
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   m.RelativeURL,
		Official:      m.Official,
	}).Validate()
}
func (m SetComicExternal) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.RID != nil {
		if *m.RID == "" {
			return GenericError("rid cannot be empty")
		}

		if len(*m.RID) != ComicGenericRIDLength {
			length := strconv.Itoa(ComicGenericRIDLength)
			return GenericError("rid must be " + length + " characters long")
		}
	}

	if err := (SetWebsite{Domain: m.WebsiteDomain}).Validate(); err != nil {
		return GenericError("website " + err.Error())
	}

	if m.RelativeURL != nil {
		if *m.RelativeURL == "" {
			return GenericError("relative url cannot be empty")
		}

		if len(*m.RelativeURL) > ComicExternalRelativeURLMax {
			max := strconv.FormatInt(ComicExternalRelativeURLMax, 10)
			return GenericError("relative url must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicSynopsisSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func (m AddComicCategory) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.CategoryID == nil && m.CategoryCode == nil {
		return GenericError("either category id or category code must exist")
	}

	return (&SetComicCategory{
		ComicID:          m.ComicID,
		ComicCode:        m.ComicCode,
		CategoryID:       m.CategoryID,
		CategoryTypeID:   m.CategoryTypeID,
		CategoryTypeCode: m.CategoryTypeCode,
		CategoryCode:     m.CategoryCode,
	}).Validate()
}
func (m SetComicCategory) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.CategoryCode != nil {
		if m.CategoryTypeID == nil && m.CategoryTypeCode == nil {
			return GenericError("either category type id or category type code must exist")
		}

		if err := (SetCategory{TypeCode: m.CategoryTypeCode, Code: m.CategoryCode}).Validate(); err != nil {
			return GenericError("category " + err.Error())
		}
	} else {
		if m.CategoryTypeID != nil || m.CategoryTypeCode != nil {
			return GenericError("category code must also be provided")
		}
	}

	return nil
}

func (m AddComicTag) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.TagID == nil && m.TagCode == nil {
		return GenericError("either tag id or tag code must exist")
	}

	return (&SetComicTag{
		ComicID:     m.ComicID,
		ComicCode:   m.ComicCode,
		TagID:       m.TagID,
		TagTypeID:   m.TagTypeID,
		TagTypeCode: m.TagTypeCode,
		TagCode:     m.TagCode,
	}).Validate()
}
func (m SetComicTag) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.TagCode != nil {
		if m.TagTypeID == nil && m.TagTypeCode == nil {
			return GenericError("either tag type id or tag type code must exist")
		}

		if err := (SetTag{TypeCode: m.TagTypeCode, Code: m.TagCode}).Validate(); err != nil {
			return GenericError("tag " + err.Error())
		}
	} else {
		if m.TagTypeID != nil || m.TagTypeCode != nil {
			return GenericError("tag code must also be provided")
		}
	}

	return nil
}

func (m AddComicRelationType) Validate() error {
	return (SetComicRelationType{
		Code: &m.Code,
		Name: &m.Name,
	}).Validate()
}
func (m SetComicRelationType) Validate() error {
	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) > ComicRelationTypeCodeMax {
			max := strconv.FormatInt(ComicRelationTypeCodeMax, 10)
			return GenericError("code must be at most " + max + " characters long")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > ComicRelationTypeNameMax {
			max := strconv.FormatInt(ComicRelationTypeNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

func (m AddComicRelation) Validate() error {
	if m.TypeID == nil && m.TypeCode == nil {
		return GenericError("either comic relation type id or comic relation type code must exist")
	}

	if m.ParentID == nil && m.ParentCode == nil {
		return GenericError("either parent comic id or parent comic code must exist")
	}

	if m.ChildID == nil && m.ChildCode == nil {
		return GenericError("either child comic id or child comic code must exist")
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
func (m SetComicRelation) Validate() error {
	if err := (SetComicRelationType{Code: m.TypeCode}).Validate(); err != nil {
		return GenericError("relation type " + err.Error())
	}

	if err := (SetComic{Code: m.ParentCode}).Validate(); err != nil {
		return GenericError("parent comic " + err.Error())
	}

	if err := (SetComic{Code: m.ChildCode}).Validate(); err != nil {
		return GenericError("child comic " + err.Error())
	}

	return nil
}

func init() {
	ComicChapterOrderByAllow = append(ComicChapterOrderByAllow, GenericOrderByAllow...)
}

const (
	ComicChapterChapterMax    = 64
	ComicChapterVersionMax    = 32
	ComicChapterVolumeMax     = 24
	ComicChapterOrderBysMax   = 5
	ComicChapterPaginationDef = 10
	ComicChapterPaginationMax = 50
	DBComicChapter            = donoengine.ID + "." + "comic_chapter"
	DBComicChapterChapter     = "chapter"
	DBComicChapterVersion     = "version"
	DBComicChapterVolume      = "volume"
	DBComicChapterReleasedAt  = "released_at"
)

var (
	ComicChapterOrderByAllow = []string{
		DBComicGenericComicID,
		DBComicChapterChapter,
		DBComicChapterVersion,
		DBComicChapterVolume,
		DBComicChapterReleasedAt,
	}

	ComicChapterSetNullAllow = []string{
		DBComicChapterChapter,
		DBComicChapterVersion,
		DBComicChapterVolume,
	}

	DBComicChapterSIDToID = func(sid ComicChapterSID) DBQueryValue {
		var comicID any
		switch {
		case sid.ComicID != nil:
			comicID = sid.ComicID
		case sid.ComicCode != nil:
			comicID = DBComicCodeToID(*sid.ComicCode)
		}
		var version any
		switch {
		case sid.Version != nil:
			version = sid.Version
		default:
			version = DBIsNull{}
		}
		return DBQueryValue{
			Table:      DBComicChapter,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: map[string]any{
				DBComicGenericComicID: comicID,
				DBComicChapterChapter: sid.Chapter,
				DBComicChapterVersion: version,
			},
		}
	}
)

type (
	ComicChapter struct {
		ID         uint       `json:"id"`
		ComicID    uint       `json:"-"`
		Chapter    string     `json:"chapter"`
		Version    *string    `json:"version"`
		Volume     *string    `json:"volume"`
		ReleasedAt time.Time  `json:"releasedAt"`
		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  *time.Time `json:"updatedAt"`
	}

	AddComicChapter struct {
		ComicID    *uint
		ComicCode  *string
		Chapter    string
		Version    *string
		Volume     *string
		ReleasedAt time.Time
	}

	SetComicChapter struct {
		ComicID    *uint
		ComicCode  *string
		Chapter    *string
		Version    *string
		Volume     *string
		ReleasedAt *time.Time
		SetNull    []string
	}

	ComicChapterSID struct {
		ComicID   *uint
		ComicCode *string
		Chapter   string
		Version   *string
	}
)

func (m AddComicChapter) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	return (SetComicChapter{
		ComicID:    m.ComicID,
		ComicCode:  m.ComicCode,
		Chapter:    &m.Chapter,
		Version:    m.Version,
		Volume:     m.Volume,
		ReleasedAt: &m.ReleasedAt,
	}).Validate()
}

func (m SetComicChapter) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.Chapter != nil {
		if *m.Chapter == "" {
			return GenericError("chapter cannot be empty")
		}

		if len(*m.Chapter) > ComicChapterChapterMax {
			max := strconv.FormatInt(ComicChapterChapterMax, 10)
			return GenericError("chapter must be at most " + max + " characters long")
		}
	}

	if m.Version != nil {
		if *m.Version == "" {
			return GenericError("version cannot be empty")
		}

		if len(*m.Version) > ComicChapterVersionMax {
			max := strconv.FormatInt(ComicChapterVersionMax, 10)
			return GenericError("version must be at most " + max + " characters long")
		}
	}

	if m.Volume != nil {
		if *m.Volume == "" {
			return GenericError("volume cannot be empty")
		}

		if len(*m.Volume) > ComicChapterVolumeMax {
			max := strconv.FormatInt(ComicChapterVolumeMax, 10)
			return GenericError("volume must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicChapterSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}
