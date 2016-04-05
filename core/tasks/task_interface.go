package tasks
import . "github.com/medcl/gopa/core/config"

type InnerTaskConfig struct {
	RuntimeConfig *RuntimeConfig
	MessageChan   *chan []byte
	QuitChan      *chan bool
	Parameter     *RoutingParameter
}

type TaskInterface interface {
	Init(config *InnerTaskConfig)
	Start() error
	Stop() error
}

