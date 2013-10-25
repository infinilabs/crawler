/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-17
 * Time: 下午5:21
 */
package types


type PathConfig struct {
	Home string
	Data string
	TaskData string
	WebData string
	Log string

	SavedFileLog string //path of saved files
	PendingFetchLog string //path of pending fetch
	FetchFailedLog string //path of failed fetch
}

type ClusterConfig struct {
	Name string
}

type RuntimeConfig struct{

	//cluster
	ClusterConfig *ClusterConfig

	//task
	TaskConfig *TaskConfig

	//splitter of joined array string
	ArrayStringSplitter string

	PathConfig *PathConfig

	GoProfEnabled bool
	StoreWebPageTogether bool

	MaxGoRoutine int

	//switch config
	ParseUrlsFromSavedPage bool
	LoadTemplatedFetchJob bool
	FetchUrlsFromSavedPage bool //fetch url parse and extracted from saved page,load data from:"pending_fetch.urls"
	ParseUrlsFromPreviousSavedPage bool //extract urls from previous saved page

	//runtime variables
	Storage Store

	WalkBloomFilterFileName string
	FetchBloomFilterFileName string
	ParseBloomFilterFileName string
	PendingFetchBloomFilterFileName string

}

//type TaskChan struct {
//	PendingFetchChan *chan []byte
//}

