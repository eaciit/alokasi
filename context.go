package alokasi

import (
    "github.com/eaciit/toolkit"
    "sync"
)

type Context struct{
    Data interface{}
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
