/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-17
 * Time: 下午12:14
 */
package store

type Store interface {
	Store(url string, data []byte)
	Get(key string) []byte
	List(from int, size int) [][]byte
	TaskEnqueue([]byte)
}

//func (b *store) Store(any interface{}){
//	return any.(store).get()
//}


