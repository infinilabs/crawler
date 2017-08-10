package model

import (
	"github.com/infinitbyte/gopa/core/store"
	"time"
)

// Domain is domain host struct
type Domain struct {
	Host       string         `storm:"id,unique" json:"host,omitempty" gorm:"not null;unique;primary_key"`
	LinksCount int64          `json:"links_count,omitempty"`
	Favicon    string         `json:"favicon,omitempty"`
	Settings   *DomainSetting `storm:"inline" json:"settings,omitempty"`
	CreateTime *time.Time     `storm:"index" json:"created,omitempty"`
	UpdateTime *time.Time     `storm:"index" json:"updated,omitempty"`
}

// DomainSetting is a settings for specific domain
type DomainSetting struct {
}

// CreateDomain create a domain
func CreateDomain(host string) Domain {
	domain := Domain{}
	domain.Host = host
	time := time.Now().UTC()
	domain.CreateTime = &time
	domain.UpdateTime = &time
	store.Create(&domain)
	return domain
}

// IncrementDomainLinkCount update domain's link count
func IncrementDomainLinkCount(host string) error {
	domain := Domain{}
	domain.Host = host

	store.Get(&domain)

	if domain.CreateTime == nil {
		domain = CreateDomain(host)
	}

	domain.LinksCount++
	store.Update(domain)

	return nil
}

// GetDomainList return domain list
func GetDomainList(from, size int, domain string) (int, []Domain, error) {
	var domains []Domain

	query := store.Query{From: from, Size: size}
	if len(domain) > 0 {
		query.Conds = store.And(store.Eq("host", domain))
	}

	err, r := store.Search(&domains, &query)

	return r.Total, domains, err
}

// GetDomain return a single domain
func GetDomain(domain string) (Domain, error) {
	var d = Domain{Host: domain}
	err := store.Get(&d)
	return d, err
}
