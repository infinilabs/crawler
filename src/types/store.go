/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-17
 * Time: 下午12:14
 */
package types


type Store interface {
	Store(url string, data []byte)
	Get(key string) []byte
	List(from int, size int) [][]byte
	TaskEnqueue([]byte)
	InitWalkBloomFilter(fileName string)
	InitFetchBloomFilter(fileName string)
	PersistBloomFilter()
	CheckWalkedUrl(url []byte) bool
	CheckFetchedUrl(url []byte) bool
	AddWalkedUrl(url []byte )
	AddFetchedUrl(url []byte )

	AddFetchFailedUrl(url []byte )

	CheckSavedFile(file string)  bool
}

//func (b *store) Store(any interface{}){
//	return any.(store).get()
//}


