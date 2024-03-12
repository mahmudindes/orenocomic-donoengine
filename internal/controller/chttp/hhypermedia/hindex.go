package hhypermedia

import (
	"context"
	"net/http"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
)

func (hpmd Hypermedia) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := hpmd.tIndex.Execute(w, struct {
		Name string
	}{
		Name: donoengine.Project,
	}); err != nil {
		log := hpmd.logger.WithContext(r.Context())
		log.ErrMessage(err, "Execute index template failed.")
	}
}

func (hpmd Hypermedia) Error(ctx context.Context, w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	if err := hpmd.tError.Execute(w, struct {
		Name       string
		StatusCode int
		Error      string
	}{
		Name:       donoengine.Project,
		StatusCode: code,
		Error:      error,
	}); err != nil {
		log := hpmd.logger.WithContext(ctx)
		log.ErrMessage(err, "Execute error template failed.")
	}
}
