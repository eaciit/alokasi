package alokasi

import (
	"github.com/eaciit/toolkit"
	"sync"
    "time"
	//"errors"
)

type AllocationTypeEnum int

const (
	AllocateAsPool AllocationTypeEnum = 1
	AllocateAsScan                    = 2
)

type Allocator struct {
	sync.RWMutex
	AllocationType AllocationTypeEnum
	WorkerNum      int
	Data           *toolkit.M

	OnRequest       func(ac *Context)
	OnReceive       func(ac *Context)
	OnSentComplete  func(ac *Context)
	OnFullyComplete func(ac *Context)

	//data interface{}
	workers               []*Worker
	wg                    *sync.WaitGroup
	sendComplete          bool
	requestingWorkerIndex int
    
    //--- channel
    cdata chan interface{}
    csend chan bool
    cdone chan bool
}

func New() *Allocator {
	a := new(Allocator)
	a.Data = &toolkit.M{}
	return a
}

func (a *Allocator) Start() {
	a.initWg()
    //a.sendComplete = make(chan bool)
	if a.WorkerNum == 0 {
		a.WorkerNum = 1
	}
	for i := 0; i < a.WorkerNum; i++ {
		w := NewWorker(a)
		w.ID = i
		a.workers = append(a.workers, w)
        w.Start()
	}
    
    a.cdata = make(chan interface{})
    a.csend = make(chan bool)
    
    go func(){
        for{
            select{
                case d:=<-a.cdata:
                    var w *Worker
                    a.initWg()
                    a.wg.Add(1)
                    w = a.requestingWorker()
                    a.setNextWorker()
                    w.Send(d)
                    break;
                    
                case <-a.csend:
                    a.sendComplete = true
                    break;
                    
                case <-time.After(1*time.Millisecond):
                    break;
            }      
        }
    }()
}

func (a *Allocator) requestingWorker() *Worker {
	if a.requestingWorkerIndex >= len(a.workers) {
		a.requestingWorkerIndex = len(a.workers) - 1
	}
	return a.workers[a.requestingWorkerIndex]
}

func (a *Allocator) setNextWorker() {
	i := a.requestingWorkerIndex
	i++
	if i >= len(a.workers) {
		i = 0
	}
	a.requestingWorkerIndex = i
}

func (a *Allocator) Send(k interface{}) {
	a.cdata <- k
}

func (a *Allocator) SendComplete() error {
	a.csend <- true
    return nil
}

func (a *Allocator) initWg() {
	if a.wg == nil {
		a.wg = new(sync.WaitGroup)
	}
}

func (a *Allocator) Wait() {
	a.initWg()
    
    for !a.sendComplete{
        time.Sleep(1*time.Millisecond)
    }
    
	a.wg.Wait()
}
