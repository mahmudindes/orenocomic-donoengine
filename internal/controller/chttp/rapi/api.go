package rapi

import (
	"context"

	"github.com/mahmudindes/orenocomic-donoengine/internal/logger"
	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

type (
	api struct {
		service Service
		oauth   OAuth
		logger  logger.Logger
	}

	Service interface {
		AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error
		GetLanguageByIETF(ctx context.Context, ietf string) (*model.Language, error)
		UpdateLanguageByIETF(ctx context.Context, ietf string, data model.SetLanguage, v *model.Language) error
		DeleteLanguageByIETF(ctx context.Context, ietf string) error
		ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error)
		CountLanguage(ctx context.Context, conds any) (int, error)

		AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error
		GetWebsiteByDomain(ctx context.Context, domain string) (*model.Website, error)
		UpdateWebsiteByDomain(ctx context.Context, domain string, data model.SetWebsite, v *model.Website) error
		DeleteWebsiteByDomain(ctx context.Context, domain string) error
		ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error)
		CountWebsite(ctx context.Context, conds any) (int, error)

		AddCategoryType(ctx context.Context, data model.AddCategoryType, v *model.CategoryType) error
		GetCategoryTypeByCode(ctx context.Context, code string) (*model.CategoryType, error)
		UpdateCategoryTypeByCode(ctx context.Context, code string, data model.SetCategoryType, v *model.CategoryType) error
		DeleteCategoryTypeByCode(ctx context.Context, code string) error
		ListCategoryType(ctx context.Context, params model.ListParams) ([]*model.CategoryType, error)
		CountCategoryType(ctx context.Context, conds any) (int, error)
		AddCategory(ctx context.Context, data model.AddCategory, v *model.Category) error
		GetCategoryBySID(ctx context.Context, sid model.CategorySID) (*model.Category, error)
		UpdateCategoryBySID(ctx context.Context, sid model.CategorySID, data model.SetCategory, v *model.Category) error
		DeleteCategoryBySID(ctx context.Context, sid model.CategorySID) error
		ListCategory(ctx context.Context, params model.ListParams) ([]*model.Category, error)
		CountCategory(ctx context.Context, conds any) (int, error)
		AddCategoryRelation(ctx context.Context, data model.AddCategoryRelation, v *model.CategoryRelation) error
		GetCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID) (*model.CategoryRelation, error)
		UpdateCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID, data model.SetCategoryRelation, v *model.CategoryRelation) error
		DeleteCategoryRelationBySID(ctx context.Context, sid model.CategoryRelationSID) error

		AddTagType(ctx context.Context, data model.AddTagType, v *model.TagType) error
		GetTagTypeByCode(ctx context.Context, code string) (*model.TagType, error)
		UpdateTagTypeByCode(ctx context.Context, code string, data model.SetTagType, v *model.TagType) error
		DeleteTagTypeByCode(ctx context.Context, code string) error
		ListTagType(ctx context.Context, params model.ListParams) ([]*model.TagType, error)
		CountTagType(ctx context.Context, conds any) (int, error)
		AddTag(ctx context.Context, data model.AddTag, v *model.Tag) error
		GetTagBySID(ctx context.Context, sid model.TagSID) (*model.Tag, error)
		UpdateTagBySID(ctx context.Context, sid model.TagSID, data model.SetTag, v *model.Tag) error
		DeleteTagBySID(ctx context.Context, sid model.TagSID) error
		ListTag(ctx context.Context, params model.ListParams) ([]*model.Tag, error)
		CountTag(ctx context.Context, conds any) (int, error)

		// Comic
		AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error
		GetComicByCode(ctx context.Context, code string) (*model.Comic, error)
		UpdateComicByCode(ctx context.Context, code string, data model.SetComic, v *model.Comic) error
		DeleteComicByCode(ctx context.Context, code string) error
		ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error)
		CountComic(ctx context.Context, conds any) (int, error)
		ExistsComicByCode(ctx context.Context, code string) (bool, error)
		AddComicTitle(ctx context.Context, data model.AddComicTitle, v *model.ComicTitle) error
		GetComicTitleBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicTitle, error)
		UpdateComicTitleBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicTitle, v *model.ComicTitle) error
		DeleteComicTitleBySID(ctx context.Context, sid model.ComicGenericSID) error
		AddComicCover(ctx context.Context, data model.AddComicCover, v *model.ComicCover) error
		GetComicCoverBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicCover, error)
		UpdateComicCoverBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicCover, v *model.ComicCover) error
		DeleteComicCoverBySID(ctx context.Context, sid model.ComicGenericSID) error
		AddComicSynopsis(ctx context.Context, data model.AddComicSynopsis, v *model.ComicSynopsis) error
		GetComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicSynopsis, error)
		UpdateComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicSynopsis, v *model.ComicSynopsis) error
		DeleteComicSynopsisBySID(ctx context.Context, sid model.ComicGenericSID) error
		AddComicExternal(ctx context.Context, data model.AddComicExternal, v *model.ComicExternal) error
		GetComicExternalBySID(ctx context.Context, sid model.ComicGenericSID) (*model.ComicExternal, error)
		UpdateComicExternalBySID(ctx context.Context, sid model.ComicGenericSID, data model.SetComicExternal, v *model.ComicExternal) error
		DeleteComicExternalBySID(ctx context.Context, sid model.ComicGenericSID) error
		AddComicCategory(ctx context.Context, data model.AddComicCategory, v *model.ComicCategory) error
		GetComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID) (*model.ComicCategory, error)
		UpdateComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID, data model.SetComicCategory, v *model.ComicCategory) error
		DeleteComicCategoryBySID(ctx context.Context, sid model.ComicCategorySID) error
		AddComicTag(ctx context.Context, data model.AddComicTag, v *model.ComicTag) error
		GetComicTagBySID(ctx context.Context, sid model.ComicTagSID) (*model.ComicTag, error)
		UpdateComicTagBySID(ctx context.Context, sid model.ComicTagSID, data model.SetComicTag, v *model.ComicTag) error
		DeleteComicTagBySID(ctx context.Context, sid model.ComicTagSID) error
		AddComicRelationType(ctx context.Context, data model.AddComicRelationType, v *model.ComicRelationType) error
		GetComicRelationTypeByCode(ctx context.Context, code string) (*model.ComicRelationType, error)
		UpdateComicRelationTypeByCode(ctx context.Context, code string, data model.SetComicRelationType, v *model.ComicRelationType) error
		DeleteComicRelationTypeByCode(ctx context.Context, code string) error
		ListComicRelationType(ctx context.Context, params model.ListParams) ([]*model.ComicRelationType, error)
		CountComicRelationType(ctx context.Context, conds any) (int, error)
		AddComicRelation(ctx context.Context, data model.AddComicRelation, v *model.ComicRelation) error
		GetComicRelationBySID(ctx context.Context, sid model.ComicRelationSID) (*model.ComicRelation, error)
		UpdateComicRelationBySID(ctx context.Context, sid model.ComicRelationSID, data model.SetComicRelation, v *model.ComicRelation) error
		DeleteComicRelationBySID(ctx context.Context, sid model.ComicRelationSID) error
		// Comic Chapter
		AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error
		GetComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (*model.ComicChapter, error)
		UpdateComicChapterBySID(ctx context.Context, sid model.ComicChapterSID, data model.SetComicChapter, v *model.ComicChapter) error
		DeleteComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) error
		ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error)
		CountComicChapter(ctx context.Context, conds any) (int, error)
	}

	OAuth interface {
		ProcessTokenContext(ctx context.Context) (bool, error)
		IsTokenExpiredError(err error) bool
	}
)

const SecuritySchemeBearerAuth = "BearerAuth"

var _ ServerInterface = (*api)(nil)

func NewAPI(svc Service, oa OAuth, log logger.Logger) *api {
	return &api{service: svc, oauth: oa, logger: log}
}
