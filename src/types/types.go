/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 上午8:56 
 */
package types

import (
	"regexp"
)

type TaskConfig struct {

	//name of this task
	Name string

	//follow page link,and walk around
	FollowLink bool

	//walking around pattern
	LinkUrlExtractRegex   *regexp.Regexp
	LinkUrlMustContain    string
	LinkUrlMustNotContain string

	//parsing url pattern,when url match this pattern,gopa will not parse urls from response of this url
	SkipPageParsePattern *regexp.Regexp

	//fetch url pattern
	FetchUrlPattern        *regexp.Regexp
	FetchUrlMustContain    string
	FetchUrlMustNotContain string

	//saving pattern
	SavingUrlPattern        *regexp.Regexp
	SavingUrlMustContain    string
	SavingUrlMustNotContain string

	//Crawling within domain
	FollowSameDomain bool
	FollowSubDomain  bool
}

type Task struct {
	Url, Request, Response []byte
}

type RoutingOffset struct {
	Partition int
	Offset uint64
}

