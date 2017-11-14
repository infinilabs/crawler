package model

import (
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
	"time"
)

// Host is host struct
type Host struct {
	Host              string          `storm:"id,unique" json:"host,omitempty" gorm:"not null;unique;primary_key" index:"id"`
	Favicon           string          `json:"favicon,omitempty"`
	CrawlerPipelineID string          `json:"crawler_pipeline_id,omitempty"`
	CrawlerPipeline   *PipelineConfig `json:"crawler_pipeline,omitempty"`
	Created           *time.Time      `storm:"index" json:"created,omitempty"`
	Updated           *time.Time      `storm:"index" json:"updated,omitempty"`
}

// CreateHost create a domain host
func CreateHost(host string) Host {
	h := Host{}
	h.Host = host
	time := time.Now().UTC()
	h.Created = &time
	h.Updated = &time
	err := persist.Save(&h)
	if err != nil {
		panic(err)
	}
	return h
}

// IncrementHostLinkCount update host's link count //TODO fix stats
func IncrementHostLinkCount(hostName string) error {
	host := Host{}
	host.Host = hostName

	persist.Get(&host)

	if host.Created == nil {
		host = CreateHost(hostName)
	}
	return nil
}

// GetHostList return host list
func GetHostList(from, size int, host string) (int, []Host, error) {
	var hosts []Host

	query := persist.Query{From: from, Size: size}
	if len(host) > 0 {
		query.Conds = persist.And(persist.Eq("host", host))
	}

	err, r := persist.Search(Host{}, &hosts, &query)

	if hosts == nil && r.Result != nil {
		t, ok := r.Result.([]interface{})
		if ok {
			for _, i := range t {
				js := util.ToJson(i, false)
				t := Host{}
				util.FromJson(js, &t)
				hosts = append(hosts, t)
			}
		}
	}

	return r.Total, hosts, err
}

// GetHost return a single host
func GetHost(host string) (Host, error) {
	var d = Host{Host: host}
	err := persist.Get(&d)
	if d.CrawlerPipelineID != "" {
		c, err := GetPipelineConfig(d.CrawlerPipelineID)
		if err != nil {
			panic(err)
		}
		d.CrawlerPipeline = c
	}
	return d, err
}
