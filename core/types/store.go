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
	Init()error
	PersistBloomFilter()
	CheckWalkedUrl(url []byte) bool
	CheckFetchedUrl(url []byte) bool
	CheckParsedFile(url []byte) bool
	AddWalkedUrl(url []byte )
	AddFetchedUrl(url []byte)
	AddSavedUrl(url []byte )   //the file already saved,but is missing in bloom filter,run this method
	AddParsedFile(url []byte )

	LogSavedFile(path,content string )

	LogFetchFailedUrl(path,content string )

//	AddFetchFailedUrl(url []byte )

	CheckSavedFile(file string)  bool

	InitPendingFetchBloomFilter(fileName string)
	CheckPendingFetchUrl(url []byte) bool
	AddPendingFetchUrl(url []byte )
	LogPendingFetchUrl(path,content string )

	LoadOffset(fileName string) int64
	PersistOffset(fileName string,offset int64)
}

//func (b *store) Store(any interface{}){
//	return any.(store).get()
//}


