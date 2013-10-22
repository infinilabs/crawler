/** 
 * User: Medcl
 * Date: 13-7-10
 * Time: 下午10:55 
 */
package fs

import (
	util "util"
	log "logging"
	."github.com/zeebo/sbloom"
	"hash/fnv"
	"io/ioutil"
	config "config"
)


type FsStore struct{
	WalkBloomFilterFileName string
	FetchBloomFilterFileName string
	WalkBloomFilter *Filter
	FetchBloomFilter *Filter
}

func (this *FsStore) Store(url string, data []byte){
	util.FilePutContentWithByte(url,data)
}

func (this *FsStore)  Get(key string) []byte {
	file,error:= util.FileGetContent(key)
	if(error!=nil){
		log.Error("get file:",key,error)
	}
	return file
}

func (this *FsStore)  List(from int, size int) [][]byte{
	  return nil
}

func (this *FsStore) TaskEnqueue(url []byte){
	 log.Info("task enqueue:",string(url))
}



func initBloomFilter(bloomFilterPersistFileName string) *Filter {
	var bloomFilter = new(Filter)
	//loading or initializing bloom filter
	if util.CheckFileExists(bloomFilterPersistFileName) {
		log.Debug("found bloomFilter,start reload,", bloomFilterPersistFileName)
		n, err := ioutil.ReadFile(bloomFilterPersistFileName)
		if err != nil {
			log.Error("bloomFilter:",bloomFilterPersistFileName, err)
		}
		if err := bloomFilter.GobDecode(n); err != nil {
			log.Error("bloomFilter:",bloomFilterPersistFileName, err)
		}
		log.Info("bloomFilter successfully reloaded:",bloomFilterPersistFileName)
	} else {
		probItems := config.GetIntConfig("BloomFilter", "ItemSize", 100000)
		log.Debug("initializing bloom-filter",bloomFilterPersistFileName,",virual size is,", probItems)
		bloomFilter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized:",bloomFilterPersistFileName)
	}
	return bloomFilter
}


func persistBloomFilter(bloomFilterPersistFileName string,bloomFilter *Filter) {

	//save bloom-filter
	m, err := bloomFilter.GobEncode()
	if err != nil {
		log.Error(err)
		return
	}
	err = ioutil.WriteFile(bloomFilterPersistFileName, m, 0600)
	if err != nil {
		panic(err)
		return
	}
	log.Info("bloomFilter safety persisted.")
}


func (this *FsStore) InitWalkBloomFilter(walkBloomFilterFileName string ){
	this.WalkBloomFilterFileName= walkBloomFilterFileName
	this.WalkBloomFilter = initBloomFilter(this.WalkBloomFilterFileName)
}

func (this *FsStore) InitFetchBloomFilter(fetchBloomFilterFileName string ){
	this.FetchBloomFilterFileName=fetchBloomFilterFileName
	this.FetchBloomFilter = initBloomFilter(this.FetchBloomFilterFileName)
}


func (this *FsStore) PersistBloomFilter(){
	persistBloomFilter(this.WalkBloomFilterFileName,this.WalkBloomFilter)
	persistBloomFilter(this.FetchBloomFilterFileName,this.FetchBloomFilter)
}

func (this *FsStore) CheckWalkedUrl(url []byte) bool{
	return this.WalkBloomFilter.Lookup(url)
}
func (this *FsStore) CheckFetchedUrl(url []byte) bool{
	return this.FetchBloomFilter.Lookup(url)
}
func (this *FsStore) AddWalkedUrl(url []byte ){
	this.WalkBloomFilter.Add(url)
}
func (this *FsStore) AddFetchedUrl(url []byte ){
	this.FetchBloomFilter.Add(url)
}

func (this *FsStore) AddFetchFailedUrl(url []byte ){
	//TODO
	log.Error("fetch failed url:",string(url))
}

func (this *FsStore) CheckSavedFile(file string)bool{
	log.Debug("start check file:",file)
	return  util.CheckFileExists(file)
}
