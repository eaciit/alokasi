package alokasi

import (
    "errors"
)

type Worker struct{
    Allocator *Allocator
    
    dataPool []interface{}
}

func NewWorker(a *Allocator) *Worker{
    w := new(Worker)
    w.Allocator = a
    return w
} 

func (w *Worker) Send(d interface{}){
   w.dataPool = append(w.dataPool, d)
}

func (w *Worker) Run()error{
    defer w.Allocator.wg.Done()
    if len(w.dataPool)==0{
        return errors.New("Worker.Run: EOF")
    }
    d := w.dataPool[0]
    ctx := NewContext(w.Allocator, d)
    w.Allocator.OnReceive(ctx)
    if ctx.Error!=nil {
        return errors.New("Worker.Run: " + ctx.Error.Error())
    }
    if len(w.dataPool)>1{
        w.dataPool = w.dataPool[1:]
    }
    return nil
}