package tests

import (
    "testing"
    "github.com/eaciit/alokasi"
    "github.com/eaciit/toolkit"
    "time"
)

func check(pre string, t *testing.T, e error){
    if e!=nil{
        t.Fatalf("%s: %s", pre, e.Error())
    }
}

/*
func TestContext(t *testing.T){
    mbase := toolkit.M{}
    c1 := alokasi.NewContext(10)
    c2 := alokasi.NewContext(20)
    c1.Items = &mbase
    c2.Items = &mbase
    mbase.Set("total",20)
    c2.Items.Set("count",4)
    toolkit.Println("C1:", toolkit.JsonString(c1))
    toolkit.Println("C2:",toolkit.JsonString(c2))
    if c1.Items.Get("total",0).(int)!=20{
        t.Fatalf("TestContext. Data is not valid")
    }
}
*/

var workerNum int = 5

func TestPool(t *testing.T){
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    alloc := alokasi.New()
    alloc.AllocationType = alokasi.AllocateAsPool
    alloc.WorkerNum = workerNum
    alloc.OnReceive = func(ac *alokasi.Context){
        ac.Allocator.Lock()
        total = ac.Allocator.Data.Get("total", 0).(int)
        workerid := ac.Setting.GetInt("workerid")
        intdata := ac.Data.(int)
        total += intdata
        toolkit.Printf("[%d] %s Processing data %d, Total now %d\n", workerid, 
            time.Now().String(),
            intdata, total)
        ac.Allocator.Data.Set("total", total)
        ac.Allocator.Unlock()
        time.Sleep(100 * time.Millisecond)
    }
    alloc.Start()
    for _, d := range data{
        go func(alloc *alokasi.Allocator, d int){
            alloc.Send(d)
        }(alloc, d)
    }
    //time.Sleep(1*time.Millisecond)
    alloc.SendComplete()   
    alloc.Wait()
    
    if total!=55{
        t.Fatalf("Total is %d, expected 55", total)
    }
}

func TestScan(t *testing.T){
    t.Skip()
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    alloc := alokasi.New()
    alloc.AllocationType = alokasi.AllocateAsScan
    alloc.WorkerNum = workerNum
    
    ikey := 0
    alloc.OnRequest = func(ac *alokasi.Context){
        //defer ac.Allocator.Unlock()
        //ac.Allocator.Lock()
        ikey = ac.Allocator.Data.GetInt("keyindex")
        if ikey==len(data){
            ac.Allocator.SendComplete()
            return
        }
        ac.Output = data[ikey]
        ikey++
        ac.Allocator.Data.Set("keyindex", ikey)
        return
    }
    alloc.OnReceive = func(ac *alokasi.Context){
        total = ac.Allocator.Data.Get("total", 0).(int)
        total += ac.Data.(int)
        ac.Allocator.Data.Set("total", total)
    }
    alloc.Start()  
    alloc.Wait()
    
    if total!=55{
        t.Fatalf("Total is %d, expected 55", total)
    }
}
