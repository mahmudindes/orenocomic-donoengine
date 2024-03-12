package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

func modelTag(m *model.Tag) Tag {
	return Tag{
		ID:        m.ID,
		TypeID:    m.TypeID,
		TypeCode:  m.TypeCode,
		Code:      m.Code,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddTag
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddTagJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag decode json body failed.")
			return
		}
		data = model.AddTag{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag parse form failed.")
			return
		}
		var data0 AddTagFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag decode form data failed.")
			return
		}
		data = model.AddTag{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	}

	result := new(model.Tag)
	if err := api.service.AddTag(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add tag failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.Code)
	response(w, modelTag(result), http.StatusCreated)
}

func (api *api) GetTag(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetTagBySID(ctx, model.TagSID{
		TypeID: &typeID,
		Code:   code,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get tag failed.")
		return
	}

	response(w, modelTag(result), http.StatusOK)
}

func (api *api) UpdateTag(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetTag
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateTagJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag decode json body failed.")
			return
		}
		data = model.SetTag{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag parse form failed.")
			return
		}
		var data0 UpdateTagFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag decode form data failed.")
			return
		}
		data = model.SetTag{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	}

	result := new(model.Tag)
	if err := api.service.UpdateTagBySID(ctx, model.TagSID{
		TypeID: &typeID,
		Code:   code,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update tag failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.Code)
	response(w, modelTag(result), http.StatusOK)
}

func (api *api) DeleteTag(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteTagBySID(ctx, model.TagSID{
		TypeID: &typeID,
		Code:   code,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete tag failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListTag(w http.ResponseWriter, r *http.Request, params ListTagParams) {
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

	totalCountCh := make(chan int, 1)
	go func() {
		count, err := api.service.CountTag(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count tag failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListTag(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List tag failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Tag
	for _, r := range result0 {
		result = append(result, modelTag(r))
	}
	response(w, result, http.StatusOK)
}
