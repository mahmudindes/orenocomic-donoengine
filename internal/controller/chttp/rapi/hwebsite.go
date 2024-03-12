package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-donoengine/internal/model"
)

func (api *api) AddWebsite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddWebsite
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddWebsiteJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add website decode json body failed.")
			return
		}
		data = model.AddWebsite{
			Domain: data0.Domain,
			Name:   data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add website parse form failed.")
			return
		}
		var data0 AddWebsiteFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add website decode form data failed.")
			return
		}
		data = model.AddWebsite{
			Domain: data0.Domain,
			Name:   data0.Name,
		}
	}

	result := new(model.Website)
	if err := api.service.AddWebsite(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add website failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Domain)
	response(w, Website{
		ID:        result.ID,
		Domain:    result.Domain,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusCreated)
}

func (api *api) GetWebsite(w http.ResponseWriter, r *http.Request, domain string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetWebsiteByDomain(ctx, domain)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get website failed.")
		return
	}

	response(w, Website{
		ID:        result.ID,
		Domain:    result.Domain,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) UpdateWebsite(w http.ResponseWriter, r *http.Request, domain string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetWebsite
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateWebsiteJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update website decode json body failed.")
			return
		}
		data = model.SetWebsite{
			Domain: data0.Domain,
			Name:   data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update website parse form failed.")
			return
		}
		var data0 UpdateWebsiteFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update website decode form data failed.")
			return
		}
		data = model.SetWebsite{
			Domain: data0.Domain,
			Name:   data0.Name,
		}
	}

	result := new(model.Website)
	if err := api.service.UpdateWebsiteByDomain(ctx, domain, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update website failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Domain)
	response(w, Website{
		ID:        result.ID,
		Domain:    result.Domain,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, http.StatusOK)
}

func (api *api) DeleteWebsite(w http.ResponseWriter, r *http.Request, domain string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteWebsiteByDomain(ctx, domain); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete website failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListWebsite(w http.ResponseWriter, r *http.Request, params ListWebsiteParams) {
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
		count, err := api.service.CountWebsite(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count website failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListWebsite(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List website failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Website
	for _, r := range result0 {
		result = append(result, Website{
			ID:        r.ID,
			Domain:    r.Domain,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	response(w, result, http.StatusOK)
}
