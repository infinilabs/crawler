package model

import (
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
	"time"
)

// Domain is domain host struct
type Domain struct {
	Host       string         `storm:"id,unique" json:"host,omitempty" gorm:"not null;unique;primary_key" index:"id"`
	LinksCount int64          `json:"links_count,omitempty"`
	Favicon    string         `json:"favicon,omitempty"`
	Settings   *DomainSetting `storm:"inline" json:"settings,omitempty"`
	Created    *time.Time     `storm:"index" json:"created,omitempty"`
	Updated    *time.Time     `storm:"index" json:"updated,omitempty"`
}

// DomainSetting is a settings for specific domain
type DomainSetting struct {
}

// CreateDomain create a domain
func CreateDomain(host string) Domain {
	domain := Domain{}
	domain.Host = host
	time := time.Now().UTC()
	domain.Created = &time
	domain.Updated = &time
	persist.Save(&domain)
	return domain
}

// IncrementDomainLinkCount update domain's link count
func IncrementDomainLinkCount(host string) error {
	domain := Domain{}
	domain.Host = host

	persist.Get(&domain)

	if domain.Created == nil {
		domain = CreateDomain(host)
	}

	domain.LinksCount++
	persist.Update(domain)

	return nil
}

// GetDomainList return domain list
func GetDomainList(from, size int, domain string) (int, []Domain, error) {
	var domains []Domain

	query := persist.Query{From: from, Size: size}
	if len(domain) > 0 {
		query.Conds = persist.And(persist.Eq("host", domain))
	}

	err, r := persist.Search(Domain{}, &domains, &query)

	if domains == nil && r.Result != nil {
		t := r.Result.([]interface{})
		for _, i := range t {
			js := util.ToJson(i, false)
			t := Domain{}
			util.FromJson(js, &t)
			domains = append(domains, t)
		}
	}

	return r.Total, domains, err
}

// GetDomain return a single domain
func GetDomain(domain string) (Domain, error) {
	var d = Domain{Host: domain}
	err := persist.Get(&d)
	return d, err
}
