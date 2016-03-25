package alokasi

import (
	"sync"
	"time"

	"github.com/eaciit/toolkit"
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
	a.cdone = make(chan bool)

	go func() {
		for {
			select {
			case d := <-a.cdata:
                go func(){
                    //a.Lock()
                    w := a.requestingWorker()
                    a.setNextWorker()
                    //a.Unlock()
                    w.Send(d)
                    //a.wg.Done()
                }()

			case <-a.csend:
				a.sendComplete = true

			case <-a.cdone:
				return

			//default:
				//-- do nothing
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
	//a.initWg()
	a.wg.Add(1)
    a.cdata <- k
    //toolkit.Println("Sending:", k)
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

	for !a.sendComplete {
		time.Sleep(1 * time.Second)
	}

	a.wg.Wait()
	a.cdone <- true
}
