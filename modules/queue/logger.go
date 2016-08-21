package queue

import logg "github.com/cihub/seelog"

type logger interface {
	Output(maxdepth int, s string) error
}

type GopaLogger struct {

}

func (log *GopaLogger) Output(maxdepth int, s string) error  {
	logg.Debug(maxdepth,s)
	return nil
}
