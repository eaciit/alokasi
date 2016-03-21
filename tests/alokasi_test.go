package tests

import (
    "testing"
    "github.com/eaciit/alokasi"
)

func TestPoolAtWorker(t *testing.T){
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    ctx := alokasi.New()
    ctx.PoolAt = alokasi.PoolAtWorker
    ctx.WorkerNum = 5
    for _, d := range data{
        ctx.Send(d)
    }
    ctx.SendComplete()   
    ctx.Wait()
    
    if total!=54{
        t.Fatalf("Total is %d, expected 54", total)
    }
}

func TestPoolAtAllocator(t *testing.T){
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    ctx := alokasi.New()
    ctx.PoolAt = alokasi.PoolAtAllocator
    ctx.WorkerNum = 5
    ctx.SetData(&data)
    ctx.Start()  
    ctx.Wait()
    
    if total!=54{
        t.Fatalf("Total is %d, expected 54", total)
    }
}
