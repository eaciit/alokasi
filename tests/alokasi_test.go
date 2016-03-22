package tests

import (
    "testing"
    "github.com/eaciit/alokasi"
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

func TestPool(t *testing.T){
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    ctx := alokasi.New()
    ctx.AllocationType = alokasi.AllocateAsPool
    ctx.WorkerNum = 5
    ctx.OnReceive = func(ac *alokasi.Context){
        ac.Lock()
        total = ac.Allocator.Data.Get("total", 0).(int)
        total += ac.Data.(int)
        ac.Allocator.Data.Set("total", total)
        ac.Unlock()
    }
    ctx.Start()
    for _, d := range data{
        ctx.Send(d)
    }
    ctx.SendComplete()   
    ctx.Wait()
    
    if total!=54{
        t.Fatalf("Total is %d, expected 54", total)
    }
}

func TestScan(t *testing.T){
    total := int(0)
    data := []int{1,2,3,4,5,6,7,8,9,10}
    
    ctx := alokasi.New()
    ctx.AllocationType = alokasi.AllocateAsScan
    ctx.WorkerNum = 5
    
    ikey := 0
    ctx.OnRequest = func(ac *alokasi.Context){
        defer ac.Unlock()
        ac.Lock()
        ikey = ac.Allocator.Data.GetInt("keyindex")
        if ikey==len(data){
            ac.SetError("EOF")
            return
        }
        ac.Output = data[ikey]
        ikey++
        ac.Allocator.Data.Set("keyindex", ikey)
        return
    }
    ctx.OnReceive = func(ac *alokasi.Context){
        ac.Lock()
        total = ac.Allocator.Data.Get("total", 0).(int)
        total += ac.Data.(int)
        ac.Allocator.Data.Set("total", total)
        ac.Unlock()
    }
    ctx.Start()  
    ctx.Wait()
    
    if total!=54{
        t.Fatalf("Total is %d, expected 54", total)
    }
}
