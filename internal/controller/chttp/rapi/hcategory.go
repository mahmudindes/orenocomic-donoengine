package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

func modelCategory(m *model.Category) Category {
	return Category{
		ID:        m.ID,
		TypeID:    m.TypeID,
		TypeCode:  m.TypeCode,
		Code:      m.Code,
		Name:      m.Name,
		Relations: slicesModel(m.Relations, modelCategoryRelation),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddCategory
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddCategoryJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category decode json body failed.")
			return
		}
		data = model.AddCategory{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category parse form failed.")
			return
		}
		var data0 AddCategoryFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category decode form data failed.")
			return
		}
		data = model.AddCategory{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	}

	result := new(model.Category)
	if err := api.service.AddCategory(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add category failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.Code)
	response(w, modelCategory(result), http.StatusCreated)
}

func (api *api) GetCategory(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetCategoryBySID(ctx, model.CategorySID{
		TypeID: &typeID,
		Code:   code,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get category failed.")
		return
	}

	response(w, modelCategory(result), http.StatusOK)
}

func (api *api) UpdateCategory(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetCategory
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateCategoryJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category decode json body failed.")
			return
		}
		data = model.SetCategory{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category parse form failed.")
			return
		}
		var data0 UpdateCategoryFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category decode form data failed.")
			return
		}
		data = model.SetCategory{
			TypeID:   data0.TypeID,
			TypeCode: data0.TypeCode,
			Code:     data0.Code,
			Name:     data0.Name,
		}
	}

	result := new(model.Category)
	if err := api.service.UpdateCategoryBySID(ctx, model.CategorySID{
		TypeID: &typeID,
		Code:   code,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update category failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+utila.Utoa(result.TypeID)+"-"+result.Code)
	response(w, modelCategory(result), http.StatusOK)
}

func (api *api) DeleteCategory(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteCategoryBySID(ctx, model.CategorySID{
		TypeID: &typeID,
		Code:   code,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete category failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListCategory(w http.ResponseWriter, r *http.Request, params ListCategoryParams) {
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
		count, err := api.service.CountCategory(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count category failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListCategory(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List category failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Category
	for _, r := range result0 {
		result = append(result, modelCategory(r))
	}
	response(w, result, http.StatusOK)
}

func modelCategoryRelation(m *model.CategoryRelation) CategoryRelation {
	return CategoryRelation{
		CategoryID:   m.ChildID,
		CategoryCode: m.ChildCode,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (api *api) AddCategoryRelation(w http.ResponseWriter, r *http.Request, typeID uint, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddCategoryRelation
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddCategoryRelationJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category relation decode json body failed.")
			return
		}
		data = model.AddCategoryRelation{
			ParentID:   nil,
			ParentCode: &code,
			TypeID:     &typeID,
			TypeCode:   nil,
			ChildID:    data0.CategoryID,
			ChildCode:  data0.CategoryCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category relation parse form failed.")
			return
		}
		var data0 AddCategoryRelationFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add category relation decode form data failed.")
			return
		}
		data = model.AddCategoryRelation{
			ParentID:   nil,
			ParentCode: &code,
			TypeID:     &typeID,
			TypeCode:   nil,
			ChildID:    data0.CategoryID,
			ChildCode:  data0.CategoryCode,
		}
	}

	result := new(model.CategoryRelation)
	if err := api.service.AddCategoryRelation(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add category relation failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.ChildCode)
	response(w, modelCategoryRelation(result), http.StatusCreated)
}

func (api *api) GetCategoryRelation(w http.ResponseWriter, r *http.Request, typeID uint, code string, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetCategoryRelationBySID(ctx, model.CategoryRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &categoryCode,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get category relation failed.")
		return
	}

	response(w, modelCategoryRelation(result), http.StatusOK)
}

func (api *api) UpdateCategoryRelation(w http.ResponseWriter, r *http.Request, typeID uint, code string, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetCategoryRelation
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateCategoryRelationJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category relation decode json body failed.")
			return
		}
		data = model.SetCategoryRelation{
			ChildID:   data0.CategoryID,
			ChildCode: data0.CategoryCode,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category relation parse form failed.")
			return
		}
		var data0 UpdateCategoryRelationFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update category relation decode form data failed.")
			return
		}
		data = model.SetCategoryRelation{
			ChildID:   data0.CategoryID,
			ChildCode: data0.CategoryCode,
		}
	}

	result := new(model.CategoryRelation)
	if err := api.service.UpdateCategoryRelationBySID(ctx, model.CategoryRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &categoryCode,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update category relation failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.ChildCode)
	response(w, modelCategoryRelation(result), http.StatusOK)
}

func (api *api) DeleteCategoryRelation(w http.ResponseWriter, r *http.Request, typeID uint, code string, categoryCode string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteCategoryRelationBySID(ctx, model.CategoryRelationSID{
		TypeID:     &typeID,
		ParentCode: &code,
		ChildCode:  &categoryCode,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete category relation failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
