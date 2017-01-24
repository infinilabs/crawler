package model

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/store"
	"time"
)

type Domain struct {
	Host       string         `storm:"id,unique" json:"host,omitempty" gorm:"not null;unique;primary_key"`
	LinksCount int64          `json:"links_count,omitempty"`
	Settings   *DomainSetting `storm:"inline" json:"settings,omitempty"`
	CreateTime *time.Time     `storm:"index" json:"created,omitempty"`
	UpdateTime *time.Time     `storm:"index" json:"updated,omitempty"`
}

type DomainSetting struct {
}

func CreateDomain(host string) Domain {
	domain := Domain{}
	domain.Host = host
	time := time.Now()
	domain.CreateTime = &time
	domain.UpdateTime = &time
	store.Create(&domain)
	return domain
}

func IncrementDomainLinkCount(host string) (error) {
	domain := Domain{}
	domain.Host = host

	store.Get(&domain)

	if domain.CreateTime == nil {
		log.Trace("create domain setting, ", host)
		domain = CreateDomain(host)
	}

	domain.LinksCount++
	store.Update(domain)

	return nil
}

func GetDomainList(from, size int, domain string) (int, []Domain, error) {
	log.Trace("start get all domain settings")
	var domains []Domain

	query := store.Query{From: from, Size: size}
	if len(domain) > 0 {
		query.Filter=&store.Cond{Name:"domain",Value:domain}
	}

	err,r:=store.Search(&domains,&query)

	return r.Total, domains, err
}
