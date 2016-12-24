/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queue

import (
	"errors"
	"github.com/medcl/gopa/core/stats"
)


type QueueKey string

type Queue interface {
	Push(QueueKey,[]byte) error
	Pop(QueueKey)([]byte)
	Close(QueueKey)(error)
}

var handler Queue

func Push(k QueueKey,v []byte) error{

	if(handler!=nil){
		o:= handler.Push(k,v)
		if(o==nil){
			stats.Increment("queue."+string(k), "push")
		}
		return o

	}
	stats.Increment("queue."+string(k), "push_error")
	panic(errors.New("channel is not registered"))
	return nil
}

func Pop(k QueueKey)([]byte){
	if(handler!=nil){
		o:=handler.Pop(k)
		stats.Increment("queue."+string(k), "pop")
		return o
	}
	stats.Increment("queue."+string(k), "pop_error")
	panic(errors.New("channel is not registered"))
}

func Close(k QueueKey)(error){
	if(handler!=nil){
		o:=handler.Close(k)
		stats.Increment("queue."+string(k), "close")
		return o
	}
	stats.Increment("queue."+string(k), "close_error")
	panic(errors.New("channel is not closed"))
}

func Register(h Queue)  {
	handler=h
}