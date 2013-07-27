/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:52 
 */
package config

type MatchRule struct {

Contain    []string
NotContain []string
Prefix []string
Suffix []string

}

type MustMatchRule struct {
	*MatchRule
}

type MustNotMatchRule struct {
	*MatchRule
}
