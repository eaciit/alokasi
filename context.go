package alokasi

import (
    "sync"
)

type Context struct{
    sync.Mutex
}