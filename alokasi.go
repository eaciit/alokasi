package alokasi

import (
    "sync"
    "github.com/eaciit/toolkit"
    //"errors"
)

type AllocationTypeEnum int

const (
	AllocateAsPool AllocationTypeEnum = 1
	AllocateAsScan               = 2
)

type Allocator struct {
	AllocationType    AllocationTypeEnum
	WorkerNum int
    Data *toolkit.M
    
    OnRequest func(ac *Context)
    OnReceive func(ac *Context)
    OnSentComplete func(ac *Context)
    OnFullyComplete func(ac *Context)
    
    //data interface{}
    workers []*Worker
    wg *sync.WaitGroup
    sendComplete bool
    requestingWorkerNum int
}

func New() *Allocator{
    a := new(Allocator)
    a.Data = &toolkit.M{}
    return a
}

func (a *Allocator) Start(){
    a.initWg()
    if a.WorkerNum==0{
        a.WorkerNum=1
    }
    for i:=0;i<a.WorkerNum;i++{
        w:=NewWorker(a)
        a.workers = append(a.workers, w)
    }
}

func (a *Allocator) requestingWorker() *Worker{
    if a.requestingWorkerNum>=len(a.workers){
        a.requestingWorkerNum=len(a.workers)-1
    }
    return a.workers[a.requestingWorkerNum]
}

func (a *Allocator) Send(k interface{}) {
    a.requestingWorker().Send(k)
}

func (a *Allocator) SendComplete() error {
	a.sendComplete = true
    return nil
}

func (a *Allocator) initWg(){
    if a.wg==nil{
        a.wg=new(sync.WaitGroup)
    }
}

func (a *Allocator) Wait() {
    a.initWg()
    a.wg.Wait()
}
