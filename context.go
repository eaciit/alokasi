package alokasi

import (
    "github.com/eaciit/toolkit"
    "sync"
    "errors"
)

type Context struct{
    Data interface{}
    Output interface{}
    Allocator *Allocator
    Setting *toolkit.M
    Error error
    
    sync.Mutex
}

func NewContext(allocator *Allocator, data interface{})*Context{
    c := new(Context)
    c.Data = data
    c.Allocator = allocator
    c.Setting = &toolkit.M{}
    return c
}

func (ac *Context) SetError(txt string){
    ac.Error = errors.New(txt)
}

func (ac *Context) Reset(){
    ac.Error = nil
}
