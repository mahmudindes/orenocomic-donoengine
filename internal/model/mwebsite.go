package model

import (
	"strconv"
	"time"

	donoengine "github.com/mahmudindes/orenocomic-donoengine"
	"github.com/mahmudindes/orenocomic-donoengine/internal/utila"
)

func init() {
	copy(WebsiteOrderByAllow, GenericOrderByAllow)
}

const (
	WebsiteDomainMax     = 32
	WebsiteNameMax       = 48
	WebsiteOrderBysMax   = 3
	WebsitePaginationDef = 10
	WebsitePaginationMax = 50
	DBWebsite            = donoengine.ID + "." + "website"
	DBWebsiteDomain      = "domain"
	DBWebsiteName        = "name"
)

var (
	WebsiteOrderByAllow = []string{
		DBWebsiteDomain,
		DBWebsiteName,
	}

	DBWebsiteDomainToID = func(domain string) DBQueryValue {
		return DBQueryValue{
			Table:      DBWebsite,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBWebsiteDomain, Value: domain},
		}
	}
)

type (
	Website struct {
		ID        uint       `json:"id"`
		Domain    string     `json:"domain"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	AddWebsite struct {
		Domain string
		Name   string
	}

	SetWebsite struct {
		Domain *string
		Name   *string
	}
)

func (m AddWebsite) Validate() error {
	return (SetWebsite{
		Domain: &m.Domain,
		Name:   &m.Name,
	}).Validate()
}

func (m SetWebsite) Validate() error {
	if m.Domain != nil {
		if *m.Domain == "" {
			return GenericError("domain cannot be empty")
		}

		if len(*m.Domain) > WebsiteDomainMax {
			max := strconv.FormatInt(WebsiteDomainMax, 10)
			return GenericError("domain must be at most " + max + " characters long")
		}

		if !utila.ValidDomain(*m.Domain) {
			return GenericError("domain is not valid")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > WebsiteNameMax {
			max := strconv.FormatInt(WebsiteNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	return nil
}

const (
	DBWebsiteGenericWebsiteID     = "website_id"
	DBWebsiteGenericWebsiteDomain = "website_domain"
)
