package model

import (
	"github.com/medcl/gopa/core/store"
	log "github.com/cihub/seelog"
	"time"
)

type Domain struct {
	Host       string         `storm:"id,unique" json:"host,omitempty"`
	LinksCount int64            `json:"links_count,omitempty"`
	Settings   *DomainSetting `storm:"inline" json:"settings,omitempty"`
	CreateTime    *time.Time       `storm:"index" json:"created,omitempty"`
	UpdateTime    *time.Time       `storm:"index" json:"updated,omitempty"`
}

type DomainSetting struct {
}

func GetDomain(host string) (error, Domain) {
	domain := Domain{}
	domain.Host = host

	time := time.Now()

	err := store.Get("Host", host, &domain)
	if err != nil {
		if(err.Error()=="not found"){
			log.Trace("create domain setting, ",host)
			domain.CreateTime = &time
			domain.UpdateTime = &time
			store.Save(&domain)
			return nil,domain
		}
	}
	//
	//domain.UpdateTime = &time
	//err = store.Update(&domain)
	//if err != nil {
	//	log.Error(err)
	//}

	return err, domain
}




func GetDomainList(from, size int, domain string) (int, []Domain, error) {
	log.Trace("start get all domain settings")
	var domains []Domain
	queryO := store.Query{Sort: "CreateTime", From: from, Size: size}
	if len(domain) > 0 {
		queryO.Filter = &store.Cond{Name: "Domain", Value: domain}
	}
	err, result := store.Search(&Domain{}, &domains, &queryO)
	if err != nil {
		log.Debug(err)
	}
	return result.Total, domains, err
}
