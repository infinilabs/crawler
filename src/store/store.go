/**
* User: Medcl
* Date: 13-7-10
* Time: 下午10:57
 */
package store

type store interface {
	store(url string, data []byte)
	get(key string) []byte
	list(from int, size int) [][]byte
}
