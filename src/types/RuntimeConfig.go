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


	//runtime variables
	Storage Store


	WalkBloomFilterFileName string
	FetchBloomFilterFileName string

}

