package common

import (
	conflog "github.com/iotexproject/Bumblebee/conf/log"
)

type InstanceOption struct {
	Logger          conflog.Logger
	Tasks           *TaskQueue
	OnStatusChanged func() // should call this when instance runtime interrupted
}

type InstanceOptionSetter func(o *InstanceOption)

func InstanceOptionWithLogger(l conflog.Logger) InstanceOptionSetter {
	return func(o *InstanceOption) { o.Logger = l }
}

var (
	DefaultLogger               = conflog.Std()
	DefaultInstanceOptionSetter = func(o *InstanceOption) {
		o.Logger = DefaultLogger
		o.Tasks = &TaskQueue{Ch: make(chan *Task, 100)}
	}
)
