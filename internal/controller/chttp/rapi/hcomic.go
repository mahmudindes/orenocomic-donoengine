package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

//
// Comic
//

func modelComic(m *model.Comic) Comic {
	return Comic{
		ID:            m.ID,
		Code:          m.Code,
		LanguageID:    m.LanguageID,
		LanguageIETF:  m.LanguageIETF,
		Titles:        slicesModel(m.Titles, modelComicTitle),
		Covers:        slicesModel(m.Covers, modelComicCover),
		Synopses:      slicesModel(m.Synopses, modelComicSynopsis),
		PublishedFrom: m.PublishedFrom,
		PublishedTo:   m.PublishedTo,
		TotalChapter:  m.TotalChapter,
		TotalVolume:   m.TotalVolume,
		NSFW:          m.NSFW,
		NSFL:          m.NSFL,
		Chapters:      slicesModel(m.Chapters, modelComicChapter),
		Externals:     slicesModel(m.Externals, modelComicExternal),
		Categories:    slicesModel(m.Categories, modelCategory),
		Tags:          slicesModel(m.Tags, modelTag),
		Relations:     slicesModel(m.Relations, modelComicRelation),
		Additionals:   m.Additionals,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (api *api) AddComic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComic
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic decode json body failed.")
			return
		}
		data = model.AddComic{
			Code:          data0.Code,
			LanguageID:    data0.LanguageID,
			LanguageIETF:  data0.LanguageIETF,
			PublishedFrom: data0.PublishedFrom,
			PublishedTo:   data0.PublishedTo,
			TotalChapter:  data0.TotalChapter,
			TotalVolume:   data0.TotalVolume,
			NSFW:          data0.NSFW,
			NSFL:          data0.NSFL,
			Additionals:   nil,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic parse form failed.")
			return
		}
		var data0 AddComicFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic decode form data failed.")
			return
		}
		data = model.AddComic{
			Code:          data0.Code,
			LanguageID:    data0.LanguageID,
			LanguageIETF:  data0.LanguageIETF,
			PublishedFrom: data0.PublishedFrom,
			PublishedTo:   data0.PublishedTo,
			TotalChapter:  data0.TotalChapter,
			TotalVolume:   data0.TotalVolume,
			NSFW:          data0.NSFW,
			NSFL:          data0.NSFL,
			Additionals:   nil,
		}
	}

	result := new(model.Comic)
	if err := api.service.AddComic(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, modelComic(result), http.StatusCreated)
}

func (api *api) GetComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicByCode(ctx, code)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic failed.")
		return
	}

	response(w, modelComic(result), http.StatusOK)
}

func (api *api) UpdateComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComic
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic decode json body failed.")
			return
		}
		data = model.SetComic{
			Code:          data0.Code,
			LanguageID:    data0.LanguageID,
			LanguageIETF:  data0.LanguageIETF,
			PublishedFrom: data0.PublishedFrom,
			PublishedTo:   data0.PublishedTo,
			TotalChapter:  data0.TotalChapter,
			TotalVolume:   data0.TotalVolume,
			NSFW:          data0.NSFW,
			NSFL:          data0.NSFL,
			Additionals:   nil,
			SetNull:       data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic parse form failed.")
			return
		}
		var data0 UpdateComicFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic decode form data failed.")
			return
		}
		data = model.SetComic{
			Code:          data0.Code,
			LanguageID:    data0.LanguageID,
			LanguageIETF:  data0.LanguageIETF,
			PublishedFrom: data0.PublishedFrom,
			PublishedTo:   data0.PublishedTo,
			TotalChapter:  data0.TotalChapter,
			TotalVolume:   data0.TotalVolume,
			NSFW:          data0.NSFW,
			NSFL:          data0.NSFL,
			Additionals:   nil,
			SetNull:       data0.SetNull,
		}
	}

	result := new(model.Comic)
	if err := api.service.UpdateComicByCode(ctx, code, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, modelComic(result), http.StatusOK)
}

func (api *api) DeleteComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicByCode(ctx, code); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListComic(w http.ResponseWriter, r *http.Request, params ListComicParams) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	pagination := model.Pagination{Page: 1, Limit: 10}
	if params.Page != nil {
		pagination.Page = *params.Page
	}
	if params.Limit != nil {
		pagination.Limit = *params.Limit
	}

	var orderBys model.OrderBys
	if params.OrderBy != nil {
		orderBys = queryOrderBys(*params.OrderBy)
	}

	conditions := []any{}

	if params.ComicExternal != nil {
		conditions1 := []any{}
		for _, ce := range *params.ComicExternal {
			conditions1a := map[string]any{}
			for _, kv := range strings.Split(ce, ",") {
				if key, val, ok := strings.Cut(kv, "="); ok {
					switch key {
					case "websiteID":
						key = model.DBWebsiteGenericWebsiteID
					case "websiteDomain":
						key = model.DBWebsiteGenericWebsiteDomain
					case "relativeURL":
						key = model.DBComicExternalRelativeURL
					case "official":
						key = model.DBComicExternalOfficial
					default:
						continue
					}
					val, err := url.QueryUnescape(val)
					if err != nil {
						continue
					}
					conditions1a[key] = val
				}
			}
			if len(conditions1a) > 0 {
				conditions1 = append(conditions1, conditions1a)
			}
		}
		if len(conditions1) > 0 {
			conditions = append(conditions, model.DBCrossConditional{
				Table:      model.DBComicExternal,
				Conditions: conditions1,
			})
		}
	}

	totalCountCh := make(chan int, 1)
	go func() {
		count, err := api.service.CountComic(ctx, conditions)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count comic failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListComic(ctx, model.ListParams{
		Conditions: conditions,
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List comic failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	result := []Comic{}
	for _, r := range result0 {
		result = append(result, modelComic(r))
	}
	response(w, result, http.StatusOK)
}

// Comic Title

func modelComicTitle(m *model.ComicTitle) ComicTitle {
	return ComicTitle{
		ID:           m.ID,
		RID:          m.RID,
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
		Title:        m.Title,
		Synonym:      m.Synonym,
		Romanized:    m.Romanized,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (api *api) AddComicTitle(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicTitle
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicTitleJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic title decode json body failed.")
			return
		}
		data = model.AddComicTitle{
			ComicID:      nil,
			ComicCode:    &code,
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Title:        data0.Title,
			Synonym:      data0.Synonym,
			Romanized:    data0.Romanized,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic title parse form failed.")
			return
		}
		var data0 AddComicTitleFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic title decode form data failed.")
			return
		}
		data = model.AddComicTitle{
			ComicID:      nil,
			ComicCode:    &code,
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Title:        data0.Title,
			Synonym:      data0.Synonym,
			Romanized:    data0.Romanized,
		}
	}

	result := new(model.ComicTitle)
	if err := api.service.AddComicTitle(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic title failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicTitle(result), http.StatusCreated)
}

func (api *api) GetComicTitle(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicTitleBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic title failed.")
		return
	}

	response(w, modelComicTitle(result), http.StatusOK)
}

func (api *api) UpdateComicTitle(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicTitle
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicTitleJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic title decode json body failed.")
			return
		}
		data = model.SetComicTitle{
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Title:        data0.Title,
			Synonym:      data0.Synonym,
			Romanized:    data0.Romanized,
			SetNull:      data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic title parse form failed.")
			return
		}
		var data0 UpdateComicTitleFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic title decode form data failed.")
			return
		}
		data = model.SetComicTitle{
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Title:        data0.Title,
			Synonym:      data0.Synonym,
			Romanized:    data0.Romanized,
			SetNull:      data0.SetNull,
		}
	}

	result := new(model.ComicTitle)
	if err := api.service.UpdateComicTitleBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic title failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicTitle(result), http.StatusOK)
}

func (api *api) DeleteComicTitle(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicTitleBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic title failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic Cover

func modelComicCover(m *model.ComicCover) ComicCover {
	return ComicCover{
		ID:            m.ID,
		RID:           m.RID,
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   m.RelativeURL,
		Priority:      m.Priority,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (api *api) AddComicCover(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicCover
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicCoverJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic cover decode json body failed.")
			return
		}
		data = model.AddComicCover{
			ComicID:       nil,
			ComicCode:     &code,
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Priority:      data0.Priority,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic cover parse form failed.")
			return
		}
		var data0 AddComicCoverFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic cover decode form data failed.")
			return
		}
		data = model.AddComicCover{
			ComicID:       nil,
			ComicCode:     &code,
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Priority:      data0.Priority,
		}
	}

	result := new(model.ComicCover)
	if err := api.service.AddComicCover(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic cover failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicCover(result), http.StatusCreated)
}

func (api *api) GetComicCover(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicCoverBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic cover failed.")
		return
	}

	response(w, modelComicCover(result), http.StatusOK)
}

func (api *api) UpdateComicCover(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicCover
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicCoverJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic cover decode json body failed.")
			return
		}
		data = model.SetComicCover{
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Priority:      data0.Priority,
			SetNull:       data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic cover parse form failed.")
			return
		}
		var data0 UpdateComicCoverFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic cover decode form data failed.")
			return
		}
		data = model.SetComicCover{
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Priority:      data0.Priority,
			SetNull:       data0.SetNull,
		}
	}

	result := new(model.ComicCover)
	if err := api.service.UpdateComicCoverBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic cover failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicCover(result), http.StatusCreated)
}

func (api *api) DeleteComicCover(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicCoverBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic cover failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic Synopsis

func modelComicSynopsis(m *model.ComicSynopsis) ComicSynopsis {
	return ComicSynopsis{
		ID:           m.ID,
		RID:          m.RID,
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
		Synopsis:     m.Synopsis,
		Version:      m.Version,
		Romanized:    m.Romanized,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (api *api) AddComicSynopsis(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicSynopsis
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicSynopsisJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic synopsis decode json body failed.")
			return
		}
		data = model.AddComicSynopsis{
			ComicID:      nil,
			ComicCode:    &code,
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Synopsis:     data0.Synopsis,
			Version:      data0.Version,
			Romanized:    data0.Romanized,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic synopsis parse form failed.")
			return
		}
		var data0 AddComicSynopsisFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic synopsis decode form data failed.")
			return
		}
		data = model.AddComicSynopsis{
			ComicID:      nil,
			ComicCode:    &code,
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Synopsis:     data0.Synopsis,
			Version:      data0.Version,
			Romanized:    data0.Romanized,
		}
	}

	result := new(model.ComicSynopsis)
	if err := api.service.AddComicSynopsis(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic synopsis failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicSynopsis(result), http.StatusCreated)
}

func (api *api) GetComicSynopsis(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicSynopsisBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic synopsis failed.")
		return
	}

	response(w, modelComicSynopsis(result), http.StatusOK)
}

func (api *api) UpdateComicSynopsis(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicSynopsis
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicSynopsisJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic synopsis decode json body failed.")
			return
		}
		data = model.SetComicSynopsis{
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Synopsis:     data0.Synopsis,
			Version:      data0.Version,
			Romanized:    data0.Romanized,
			SetNull:      data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic synopsis parse form failed.")
			return
		}
		var data0 UpdateComicSynopsisFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic synopsis decode form data failed.")
			return
		}
		data = model.SetComicSynopsis{
			RID:          data0.RID,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
			Synopsis:     data0.Synopsis,
			Version:      data0.Version,
			Romanized:    data0.Romanized,
			SetNull:      data0.SetNull,
		}
	}

	result := new(model.ComicSynopsis)
	if err := api.service.UpdateComicSynopsisBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic synopsis failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicSynopsis(result), http.StatusOK)
}

func (api *api) DeleteComicSynopsis(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicSynopsisBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic synopsis failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic External

func modelComicExternal(m *model.ComicExternal) ComicExternal {
	return ComicExternal{
		ID:            m.ID,
		RID:           m.RID,
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   m.RelativeURL,
		Official:      m.Official,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (api *api) AddComicExternal(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicExternal
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicExternalJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic external decode json body failed.")
			return
		}
		data = model.AddComicExternal{
			ComicID:       nil,
			ComicCode:     &code,
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Official:      data0.Official,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic external parse form failed.")
			return
		}
		var data0 AddComicExternalFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic external decode form data failed.")
			return
		}
		data = model.AddComicExternal{
			ComicID:       nil,
			ComicCode:     &code,
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Official:      data0.Official,
		}
	}

	result := new(model.ComicExternal)
	if err := api.service.AddComicExternal(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic external failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicExternal(result), http.StatusCreated)
}

func (api *api) GetComicExternal(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicExternalBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic external failed.")
		return
	}

	response(w, modelComicExternal(result), http.StatusOK)
}

func (api *api) UpdateComicExternal(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicExternal
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicExternalJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic external decode json body failed.")
			return
		}
		data = model.SetComicExternal{
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Official:      data0.Official,
			SetNull:       data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic external parse form failed.")
			return
		}
		var data0 UpdateComicExternalFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic external decode form data failed.")
			return
		}
		data = model.SetComicExternal{
			RID:           data0.RID,
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			Official:      data0.Official,
			SetNull:       data0.SetNull,
		}
	}

	result := new(model.ComicExternal)
	if err := api.service.UpdateComicExternalBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic external failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.RID)
	response(w, modelComicExternal(result), http.StatusOK)
}

func (api *api) DeleteComicExternal(w http.ResponseWriter, r *http.Request, code string, rid string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicExternalBySID(ctx, model.ComicGenericSID{
		ComicCode: &code,
		RID:       rid,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic external failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic Category

func modelComicCategory(m *model.ComicCategory) ComicCategory {
	return ComicCategory{
		CategoryID:     m.CategoryID,
		CategoryTypeID: m.CategoryTypeID,
		CategoryCode:   m.CategoryCode,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (api *api) AddComicCategory(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicCategory
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicCategoryJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic category decode json body failed.")
			return
		}
		data = model.AddComicCategory{
			ComicID:          nil,
			ComicCode:        &code,
			CategoryID:       data0.CategoryID,
			CategoryTypeID:   data0.CategoryTypeID,
			CategoryTypeCode: data0.CategoryTypeCode,
			CategoryCode:     data0.CategoryCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic category parse form failed.")
			return
		}
		var data0 AddComicCategoryFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic category decode form data failed.")
			return
		}
		data = model.AddComicCategory{
			ComicID:          nil,
			ComicCode:        &code,
			CategoryID:       data0.CategoryID,
			CategoryTypeID:   data0.CategoryTypeID,
			CategoryTypeCode: data0.CategoryTypeCode,
			CategoryCode:     data0.CategoryCode,
		}
	}

	result := new(model.ComicCategory)
	if err := api.service.AddComicCategory(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic category failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.CategoryTypeID)+"-"+result.CategoryCode)
	response(w, modelComicCategory(result), http.StatusCreated)
}

func (api *api) GetComicCategory(w http.ResponseWriter, r *http.Request, code string, typeID uint, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicCategoryBySID(ctx, model.ComicCategorySID{
		ComicCode:   &code,
		CategorySID: &model.CategorySID{TypeID: &typeID, Code: categoryCode},
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic category failed.")
		return
	}

	response(w, modelComicCategory(result), http.StatusOK)
}

func (api *api) UpdateComicCategory(w http.ResponseWriter, r *http.Request, code string, typeID uint, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicCategory
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicCategoryJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic category decode json body failed.")
			return
		}
		data = model.SetComicCategory{
			CategoryID:       data0.CategoryID,
			CategoryTypeID:   data0.CategoryTypeID,
			CategoryTypeCode: data0.CategoryTypeCode,
			CategoryCode:     data0.CategoryCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic category parse form failed.")
			return
		}
		var data0 UpdateComicCategoryFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic category decode form data failed.")
			return
		}
		data = model.SetComicCategory{
			CategoryID:       data0.CategoryID,
			CategoryTypeID:   data0.CategoryTypeID,
			CategoryTypeCode: data0.CategoryTypeCode,
			CategoryCode:     data0.CategoryCode,
		}
	}

	result := new(model.ComicCategory)
	if err := api.service.UpdateComicCategoryBySID(ctx, model.ComicCategorySID{
		ComicCode:   &code,
		CategorySID: &model.CategorySID{TypeID: &typeID, Code: categoryCode},
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic category failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.CategoryTypeID)+"-"+result.CategoryCode)
	response(w, modelComicCategory(result), http.StatusOK)
}

func (api *api) DeleteComicCategory(w http.ResponseWriter, r *http.Request, code string, typeID uint, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicCategoryBySID(ctx, model.ComicCategorySID{
		ComicCode:   &code,
		CategorySID: &model.CategorySID{TypeID: &typeID, Code: categoryCode},
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic category failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic Tag

func modelComicTag(m *model.ComicTag) ComicTag {
	return ComicTag{
		TagID:     m.TagID,
		TagTypeID: m.TagTypeID,
		TagCode:   m.TagCode,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddComicTag(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicTag
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicTagJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic tag decode json body failed.")
			return
		}
		data = model.AddComicTag{
			ComicID:     nil,
			ComicCode:   &code,
			TagID:       data0.TagID,
			TagTypeID:   data0.TagTypeID,
			TagTypeCode: data0.TagTypeCode,
			TagCode:     data0.TagCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic tag parse form failed.")
			return
		}
		var data0 AddComicTagFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic tag decode form data failed.")
			return
		}
		data = model.AddComicTag{
			ComicID:     nil,
			ComicCode:   &code,
			TagID:       data0.TagID,
			TagTypeID:   data0.TagTypeID,
			TagTypeCode: data0.TagTypeCode,
			TagCode:     data0.TagCode,
		}
	}

	result := new(model.ComicTag)
	if err := api.service.AddComicTag(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic tag failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TagTypeID)+"-"+result.TagCode)
	response(w, modelComicTag(result), http.StatusCreated)
}

func (api *api) GetComicTag(w http.ResponseWriter, r *http.Request, code string, typeID uint, tagCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicTagBySID(ctx, model.ComicTagSID{
		ComicCode: &code,
		TagSID:    &model.TagSID{TypeID: &typeID, Code: tagCode},
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic tag failed.")
		return
	}

	response(w, modelComicTag(result), http.StatusOK)
}

func (api *api) UpdateComicTag(w http.ResponseWriter, r *http.Request, code string, typeID uint, tagCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicTag
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicTagJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic tag decode json body failed.")
			return
		}
		data = model.SetComicTag{
			TagID:       data0.TagID,
			TagTypeID:   data0.TagTypeID,
			TagTypeCode: data0.TagTypeCode,
			TagCode:     data0.TagCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic tag parse form failed.")
			return
		}
		var data0 UpdateComicTagFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic tag decode form data failed.")
			return
		}
		data = model.SetComicTag{
			TagID:       data0.TagID,
			TagTypeID:   data0.TagTypeID,
			TagTypeCode: data0.TagTypeCode,
			TagCode:     data0.TagCode,
		}
	}

	result := new(model.ComicTag)
	if err := api.service.UpdateComicTagBySID(ctx, model.ComicTagSID{
		ComicCode: &code,
		TagSID:    &model.TagSID{TypeID: &typeID, Code: tagCode},
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic tag failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TagTypeID)+"-"+result.TagCode)
	response(w, modelComicTag(result), http.StatusOK)
}

func (api *api) DeleteComicTag(w http.ResponseWriter, r *http.Request, code string, typeID uint, tagCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicTagBySID(ctx, model.ComicTagSID{
		ComicCode: &code,
		TagSID:    &model.TagSID{TypeID: &typeID, Code: tagCode},
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic tag failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comic Relation

func modelComicRelation(m *model.ComicRelation) ComicRelation {
	return ComicRelation{
		TypeID:    m.TypeID,
		ComicID:   m.ChildID,
		ComicCode: m.ChildCode,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddComicRelation(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicRelation
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicRelationJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation decode json body failed.")
			return
		}
		data = model.AddComicRelation{
			ParentID:   nil,
			ParentCode: &code,
			TypeID:     data0.TypeID,
			TypeCode:   data0.TypeCode,
			ChildID:    data0.ComicID,
			ChildCode:  data0.ComicCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation parse form failed.")
			return
		}
		var data0 AddComicRelationFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation decode form data failed.")
			return
		}
		data = model.AddComicRelation{
			ParentID:   nil,
			ParentCode: &code,
			TypeID:     data0.TypeID,
			TypeCode:   data0.TypeCode,
			ChildID:    data0.ComicID,
			ChildCode:  data0.ComicCode,
		}
	}

	result := new(model.ComicRelation)
	if err := api.service.AddComicRelation(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic relation failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.ChildCode)
	response(w, modelComicRelation(result), http.StatusCreated)
}

func (api *api) GetComicRelation(w http.ResponseWriter, r *http.Request, code string, typeID uint, comicCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicRelationBySID(ctx, model.ComicRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &comicCode,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic relation failed.")
		return
	}

	response(w, modelComicRelation(result), http.StatusOK)
}

func (api *api) UpdateComicRelation(w http.ResponseWriter, r *http.Request, code string, typeID uint, comicCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicRelation
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicRelationJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation decode json body failed.")
			return
		}
		data = model.SetComicRelation{
			TypeID:    data0.TypeID,
			TypeCode:  data0.TypeCode,
			ChildID:   data0.ComicID,
			ChildCode: data0.ComicCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation parse form failed.")
			return
		}
		var data0 UpdateComicRelationFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation decode form data failed.")
			return
		}
		data = model.SetComicRelation{
			TypeID:    data0.TypeID,
			TypeCode:  data0.TypeCode,
			ChildID:   data0.ComicID,
			ChildCode: data0.ComicCode,
		}
	}

	result := new(model.ComicRelation)
	if err := api.service.UpdateComicRelationBySID(ctx, model.ComicRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &comicCode,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic relation failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.ChildCode)
	response(w, modelComicRelation(result), http.StatusOK)
}

func (api *api) DeleteComicRelation(w http.ResponseWriter, r *http.Request, code string, typeID uint, comicCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicRelationBySID(ctx, model.ComicRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &comicCode,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic relation failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//
// Comic Chapter
//

func modelComicChapter(m *model.ComicChapter) ComicChapter {
	return ComicChapter{
		ID:         m.ID,
		Chapter:    m.Chapter,
		Version:    m.Version,
		Volume:     m.Volume,
		ReleasedAt: m.ReleasedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (api *api) AddComicChapter(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicChapter
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicChapterJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter decode json body failed.")
			return
		}
		data = model.AddComicChapter{
			ComicID:    nil,
			ComicCode:  &code,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			Volume:     data0.Volume,
			ReleasedAt: data0.ReleasedAt,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter parse form failed.")
			return
		}
		var data0 AddComicChapterFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter decode form data failed.")
			return
		}
		data = model.AddComicChapter{
			ComicID:    nil,
			ComicCode:  &code,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			Volume:     data0.Volume,
			ReleasedAt: data0.ReleasedAt,
		}
	}

	result := new(model.ComicChapter)
	if err := api.service.AddComicChapter(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic chapter failed.")
		return
	}

	slug := url.QueryEscape(result.Chapter)
	if result.Version != nil {
		slug += "+" + url.QueryEscape(*result.Version)
	}

	w.Header().Set("Location", r.URL.Path+"/"+slug)
	response(w, modelComicChapter(result), http.StatusCreated)
}

func (api *api) GetComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic chapter failed.")
		return
	}

	response(w, modelComicChapter(result), http.StatusOK)
}

func (api *api) UpdateComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicChapter
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicChapterJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter decode json body failed.")
			return
		}
		data = model.SetComicChapter{
			ComicID:    nil,
			ComicCode:  nil,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			Volume:     data0.Volume,
			ReleasedAt: data0.ReleasedAt,
			SetNull:    data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter parse form failed.")
			return
		}
		var data0 UpdateComicChapterFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter decode form data failed.")
			return
		}
		data = model.SetComicChapter{
			ComicID:    nil,
			ComicCode:  nil,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			Volume:     data0.Volume,
			ReleasedAt: data0.ReleasedAt,
			SetNull:    data0.SetNull,
		}
	}

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	result := new(model.ComicChapter)
	if err := api.service.UpdateComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic chapter failed.")
		return
	}

	slug := url.QueryEscape(result.Chapter)
	if result.Version != nil {
		slug += "+" + url.QueryEscape(*result.Version)
	}

	w.Header().Set("Location", r.URL.Path+"/"+slug)
	response(w, modelComicChapter(result), http.StatusOK)
}

func (api *api) DeleteComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic chapter failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListComicChapter(w http.ResponseWriter, r *http.Request, code string, params ListComicChapterParams) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	pagination := model.Pagination{Page: 1, Limit: 10}
	if params.Page != nil {
		pagination.Page = *params.Page
	}
	if params.Limit != nil {
		pagination.Limit = *params.Limit
	}

	var orderBys model.OrderBys
	if params.OrderBy != nil {
		orderBys = queryOrderBys(*params.OrderBy)
	}

	conditions := model.DBConditionalKV{
		Key:   model.DBComicGenericComicID,
		Value: model.DBComicCodeToID(code),
	}

	totalCountCh := make(chan int, 1)
	go func() {
		count, err := api.service.CountComicChapter(ctx, conditions)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count comic chapter failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListComicChapter(ctx, model.ListParams{
		Conditions: conditions,
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List comic chapter failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []ComicChapter
	for _, r := range result0 {
		result = append(result, modelComicChapter(r))
	}
	response(w, result, http.StatusOK)
}
