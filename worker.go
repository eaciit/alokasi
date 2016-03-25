package alokasi

import (
	//"time"
	"github.com/eaciit/toolkit"
	//"strings"
)

type WorkerStateEnum int

const (
	WorkerIdle    WorkerStateEnum = 0
	WorkerRunning                 = 1
	WorkerStop                    = 2
)

type Worker struct {
	ID        int
	Allocator *Allocator
	Status    WorkerStateEnum

	//dataPool []interface{}
	Setting *toolkit.M
    chanData chan interface{}
}

func NewWorker(a *Allocator) *Worker {
	w := new(Worker)
	w.Allocator = a
	w.Status = WorkerIdle
    w.Setting = &toolkit.M{}
    w.chanData = make(chan interface{})
	return w
}

func (w *Worker) Start(){
    if w.Allocator.AllocationType==AllocateAsPool{
        w.startAsPool()
    }
}

func (w *Worker) startAsPool() {
	go func() {
		for {
			select {
			case d:= <-w.chanData:
				w.exec(d)

			default:
            }
		}
	}()
}

func (w *Worker) stop() {
	w.Status = WorkerStop
}


func (w *Worker) Send(d interface{}) {
	//toolkit.Printf("Receving %v to datapool of worker %d\n", d, w.ID)
    w.chanData <- d
}

func (w *Worker) exec(d interface{}){
    defer func(){
        w.Allocator.wg.Done()
    }()
    ctx := NewContext(w, d)
    ctx.Setting.Set("workerid", w.ID)
  	w.Allocator.OnReceive(ctx)	
}
