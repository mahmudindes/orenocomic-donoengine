package service

import (
	"context"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

type (
	Service struct {
		database database
		oauth    oauth
	}

	database interface {
		AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error
		GetLanguage(ctx context.Context, conds any) (*model.Language, error)
		UpdateLanguage(ctx context.Context, data model.SetLanguage, conds any, v *model.Language) error
		DeleteLanguage(ctx context.Context, conds any, v *model.Language) error
		ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error)
		CountLanguage(ctx context.Context, conds any) (int, error)

		AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error
		GetWebsite(ctx context.Context, conds any) (*model.Website, error)
		UpdateWebsite(ctx context.Context, data model.SetWebsite, conds any, v *model.Website) error
		DeleteWebsite(ctx context.Context, conds any, v *model.Website) error
		ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error)
		CountWebsite(ctx context.Context, conds any) (int, error)

		AddCategoryType(ctx context.Context, data model.AddCategoryType, v *model.CategoryType) error
		GetCategoryType(ctx context.Context, conds any) (*model.CategoryType, error)
		UpdateCategoryType(ctx context.Context, data model.SetCategoryType, conds any, v *model.CategoryType) error
		DeleteCategoryType(ctx context.Context, conds any, v *model.CategoryType) error
		ListCategoryType(ctx context.Context, params model.ListParams) ([]*model.CategoryType, error)
		CountCategoryType(ctx context.Context, conds any) (int, error)
		AddCategory(ctx context.Context, data model.AddCategory, v *model.Category) error
		GetCategory(ctx context.Context, conds any) (*model.Category, error)
		UpdateCategory(ctx context.Context, data model.SetCategory, conds any, v *model.Category) error
		DeleteCategory(ctx context.Context, conds any, v *model.Category) error
		ListCategory(ctx context.Context, params model.ListParams) ([]*model.Category, error)
		CountCategory(ctx context.Context, conds any) (int, error)
		AddCategoryRelation(ctx context.Context, data model.AddCategoryRelation, v *model.CategoryRelation) error
		GetCategoryRelation(ctx context.Context, conds any) (*model.CategoryRelation, error)
		UpdateCategoryRelation(ctx context.Context, data model.SetCategoryRelation, conds any, v *model.CategoryRelation) error
		DeleteCategoryRelation(ctx context.Context, conds any, v *model.CategoryRelation) error
		ListCategoryRelation(ctx context.Context, params model.ListParams) ([]*model.CategoryRelation, error)
		CountCategoryRelation(ctx context.Context, conds any) (int, error)

		AddTagType(ctx context.Context, data model.AddTagType, v *model.TagType) error
		GetTagType(ctx context.Context, conds any) (*model.TagType, error)
		UpdateTagType(ctx context.Context, data model.SetTagType, conds any, v *model.TagType) error
		DeleteTagType(ctx context.Context, conds any, v *model.TagType) error
		ListTagType(ctx context.Context, params model.ListParams) ([]*model.TagType, error)
		CountTagType(ctx context.Context, conds any) (int, error)
		AddTag(ctx context.Context, data model.AddTag, v *model.Tag) error
		GetTag(ctx context.Context, conds any) (*model.Tag, error)
		UpdateTag(ctx context.Context, data model.SetTag, conds any, v *model.Tag) error
		DeleteTag(ctx context.Context, conds any, v *model.Tag) error
		ListTag(ctx context.Context, params model.ListParams) ([]*model.Tag, error)
		CountTag(ctx context.Context, conds any) (int, error)

		// Comic
		AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error
		GetComic(ctx context.Context, conds any) (*model.Comic, error)
		UpdateComic(ctx context.Context, data model.SetComic, conds any, v *model.Comic) error
		DeleteComic(ctx context.Context, conds any, v *model.Comic) error
		ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error)
		CountComic(ctx context.Context, conds any) (int, error)
		ExistsComic(ctx context.Context, conds any) (bool, error)
		AddComicTitle(ctx context.Context, data model.AddComicTitle, v *model.ComicTitle) error
		GetComicTitle(ctx context.Context, conds any) (*model.ComicTitle, error)
		UpdateComicTitle(ctx context.Context, data model.SetComicTitle, conds any, v *model.ComicTitle) error
		DeleteComicTitle(ctx context.Context, conds any, v *model.ComicTitle) error
		ListComicTitle(ctx context.Context, params model.ListParams) ([]*model.ComicTitle, error)
		CountComicTitle(ctx context.Context, conds any) (int, error)
		AddComicCover(ctx context.Context, data model.AddComicCover, v *model.ComicCover) error
		GetComicCover(ctx context.Context, conds any) (*model.ComicCover, error)
		UpdateComicCover(ctx context.Context, data model.SetComicCover, conds any, v *model.ComicCover) error
		DeleteComicCover(ctx context.Context, conds any, v *model.ComicCover) error
		ListComicCover(ctx context.Context, params model.ListParams) ([]*model.ComicCover, error)
		CountComicCover(ctx context.Context, conds any) (int, error)
		AddComicSynopsis(ctx context.Context, data model.AddComicSynopsis, v *model.ComicSynopsis) error
		GetComicSynopsis(ctx context.Context, conds any) (*model.ComicSynopsis, error)
		UpdateComicSynopsis(ctx context.Context, data model.SetComicSynopsis, conds any, v *model.ComicSynopsis) error
		DeleteComicSynopsis(ctx context.Context, conds any, v *model.ComicSynopsis) error
		ListComicSynopsis(ctx context.Context, params model.ListParams) ([]*model.ComicSynopsis, error)
		CountComicSynopsis(ctx context.Context, conds any) (int, error)
		AddComicExternal(ctx context.Context, data model.AddComicExternal, v *model.ComicExternal) error
		GetComicExternal(ctx context.Context, conds any) (*model.ComicExternal, error)
		UpdateComicExternal(ctx context.Context, data model.SetComicExternal, conds any, v *model.ComicExternal) error
		DeleteComicExternal(ctx context.Context, conds any, v *model.ComicExternal) error
		ListComicExternal(ctx context.Context, params model.ListParams) ([]*model.ComicExternal, error)
		CountComicExternal(ctx context.Context, conds any) (int, error)
		AddComicCategory(ctx context.Context, data model.AddComicCategory, v *model.ComicCategory) error
		GetComicCategory(ctx context.Context, conds any) (*model.ComicCategory, error)
		UpdateComicCategory(ctx context.Context, data model.SetComicCategory, conds any, v *model.ComicCategory) error
		DeleteComicCategory(ctx context.Context, conds any, v *model.ComicCategory) error
		ListComicCategory(ctx context.Context, params model.ListParams) ([]*model.ComicCategory, error)
		CountComicCategory(ctx context.Context, conds any) (int, error)
		AddComicTag(ctx context.Context, data model.AddComicTag, v *model.ComicTag) error
		GetComicTag(ctx context.Context, conds any) (*model.ComicTag, error)
		UpdateComicTag(ctx context.Context, data model.SetComicTag, conds any, v *model.ComicTag) error
		DeleteComicTag(ctx context.Context, conds any, v *model.ComicTag) error
		ListComicTag(ctx context.Context, params model.ListParams) ([]*model.ComicTag, error)
		CountComicTag(ctx context.Context, conds any) (int, error)
		AddComicRelationType(ctx context.Context, data model.AddComicRelationType, v *model.ComicRelationType) error
		GetComicRelationType(ctx context.Context, conds any) (*model.ComicRelationType, error)
		UpdateComicRelationType(ctx context.Context, data model.SetComicRelationType, conds any, v *model.ComicRelationType) error
		DeleteComicRelationType(ctx context.Context, conds any, v *model.ComicRelationType) error
		ListComicRelationType(ctx context.Context, params model.ListParams) ([]*model.ComicRelationType, error)
		CountComicRelationType(ctx context.Context, conds any) (int, error)
		AddComicRelation(ctx context.Context, data model.AddComicRelation, v *model.ComicRelation) error
		GetComicRelation(ctx context.Context, conds any) (*model.ComicRelation, error)
		UpdateComicRelation(ctx context.Context, data model.SetComicRelation, conds any, v *model.ComicRelation) error
		DeleteComicRelation(ctx context.Context, conds any, v *model.ComicRelation) error
		ListComicRelation(ctx context.Context, params model.ListParams) ([]*model.ComicRelation, error)
		CountComicRelation(ctx context.Context, conds any) (int, error)
		// Comic Chapter
		AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error
		GetComicChapter(ctx context.Context, conds any) (*model.ComicChapter, error)
		UpdateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error
		DeleteComicChapter(ctx context.Context, conds any, v *model.ComicChapter) error
		ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error)
		CountComicChapter(ctx context.Context, conds any) (int, error)
	}

	oauth interface {
		HasPermissionContext(ctx context.Context, permission string) bool
		TokenPermissionKey(s ...string) string
	}
)

func New(db database, oa oauth) Service {
	return Service{database: db, oauth: oa}
}
