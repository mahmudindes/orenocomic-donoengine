package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

//
// Category
//

func (api *api) AddCategoryType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddCategoryType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddCategoryTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category type decode json body failed.")
			return
		}
		data = model.AddCategoryType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category type parse form failed.")
			return
		}
		var data0 AddCategoryTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category type decode form data failed.")
			return
		}
		data = model.AddCategoryType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.CategoryType)
	if err := api.service.AddCategoryType(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add category type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusCreated)
}

func (api *api) GetCategoryType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetCategoryTypeByCode(ctx, code)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get category type failed.")
		return
	}

	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) UpdateCategoryType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetCategoryType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateCategoryTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category type decode json body failed.")
			return
		}
		data = model.SetCategoryType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category type parse form failed.")
			return
		}
		var data0 UpdateCategoryTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category type decode form data failed.")
			return
		}
		data = model.SetCategoryType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.CategoryType)
	if err := api.service.UpdateCategoryTypeByCode(ctx, code, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update category type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) DeleteCategoryType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteCategoryTypeByCode(ctx, code); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete category type failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListCategoryType(w http.ResponseWriter, r *http.Request, params ListCategoryTypeParams) {
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
		count, err := api.service.CountCategoryType(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count category type failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListCategoryType(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List category type failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []GenericType
	for _, r := range result0 {
		result = append(result, GenericType{
			ID:        r.ID,
			Code:      r.Code,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	response(w, result, http.StatusOK)
}

//
// Tag
//

func (api *api) AddTagType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddTagType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddTagTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag type decode json body failed.")
			return
		}
		data = model.AddTagType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag type parse form failed.")
			return
		}
		var data0 AddTagTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add tag type decode form data failed.")
			return
		}
		data = model.AddTagType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.TagType)
	if err := api.service.AddTagType(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add tag type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusCreated)
}

func (api *api) GetTagType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetTagTypeByCode(ctx, code)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get tag type failed.")
		return
	}

	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) UpdateTagType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetTagType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateTagTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag type decode json body failed.")
			return
		}
		data = model.SetTagType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag type parse form failed.")
			return
		}
		var data0 UpdateTagTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update tag type decode form data failed.")
			return
		}
		data = model.SetTagType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.TagType)
	if err := api.service.UpdateTagTypeByCode(ctx, code, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update tag type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) DeleteTagType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteTagTypeByCode(ctx, code); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete tag type failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListTagType(w http.ResponseWriter, r *http.Request, params ListTagTypeParams) {
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
		count, err := api.service.CountTagType(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count tag type failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListTagType(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List tag type failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []GenericType
	for _, r := range result0 {
		result = append(result, GenericType{
			ID:        r.ID,
			Code:      r.Code,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	response(w, result, http.StatusOK)
}

//
// Comic Relation
//

func (api *api) AddComicRelationType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicRelationType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicRelationTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation type decode json body failed.")
			return
		}
		data = model.AddComicRelationType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation type parse form failed.")
			return
		}
		var data0 AddComicRelationTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic relation type decode form data failed.")
			return
		}
		data = model.AddComicRelationType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.ComicRelationType)
	if err := api.service.AddComicRelationType(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic relation type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusCreated)
}

func (api *api) GetComicRelationType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicRelationTypeByCode(ctx, code)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic relation type failed.")
		return
	}

	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) UpdateComicRelationType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicRelationType
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicRelationTypeJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation type decode json body failed.")
			return
		}
		data = model.SetComicRelationType{
			Code: data0.Code,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation type parse form failed.")
			return
		}
		var data0 UpdateComicRelationTypeFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic relation type decode form data failed.")
			return
		}
		data = model.SetComicRelationType{
			Code: data0.Code,
			Name: data0.Name,
		}
	}

	result := new(model.ComicRelationType)
	if err := api.service.UpdateComicRelationTypeByCode(ctx, code, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic relation type failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, GenericType{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) DeleteComicRelationType(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicRelationTypeByCode(ctx, code); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic relation type failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListComicRelationType(w http.ResponseWriter, r *http.Request, params ListComicRelationTypeParams) {
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
		count, err := api.service.CountComicRelationType(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count comic relation type failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListComicRelationType(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List comic relation type failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []GenericType
	for _, r := range result0 {
		result = append(result, GenericType{
			ID:        r.ID,
			Code:      r.Code,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	response(w, result, http.StatusOK)
}
