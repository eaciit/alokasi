package alokasi

import (
    "sync"
    "github.com/eaciit/toolkit"
    "errors"
)

type AllocationTypeEnum int

const (
	AllocateAsPool AllocationTypeEnum = 1
	AllocateAsScan               = 2
)

type Allocator struct {
	AllocationType    AllocationTypeEnum
	WorkerNum int
    
    OnRequest interface{}
    OnReceive interface{}
    OnSentComplete interface{}
    OnFullyComplete interface{}
    
    //data interface{}
    wg *sync.WaitGroup
    sendComplete bool
}

func New() *Allocator{
    a := new(Allocator)
    return a
}

func (a *Allocator) Start(){
    a.wg = new(sync.WaitGroup)
}

func (a *Allocator) Send(k interface{}) {
}

func (a *Allocator) SendComplete() error {
	a.sendComplete = true
    return nil
}

func (a *Allocator) Wait() {
    a.wg.Wait()
}
