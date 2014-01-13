/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:52 
 */
package types

type MatchRule struct {

Contain    []string
NotContain []string
Prefix []string
Suffix []string

}

type ShouldMatchRule struct {
	*MatchRule
}

type MustMatchRule struct {
	*MatchRule
}

type MustNotMatchRule struct {
	*MatchRule
}
