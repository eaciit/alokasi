package tests

import (
    "testing"
    "github.com/eaciit/alokasi"
    "github.com/eaciit/toolkit"
    "sync"
    //"time"
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

var workerNum int = 100
var dataSample int = 1000
var expectation int
var data []int

func TestPrepare(t *testing.T){
    for i:=0;i<dataSample;i++{
        n := toolkit.RandInt(100)
        data = append(data, n)
        expectation +=n
    }
    //toolkit.Println("Data: ", data)
    toolkit.Println("Expecting total is", expectation)
}

func TestCalc(t *testing.T){
    total := 0
    for _, d := range data{
        total += d
    }
    if total!=expectation{
        t.Fatalf("Total is %d, expected %d", total, expectation)
    }
}

func TestSimpleChannel(t *testing.T){
    total := 0
    processed := 0
    cdata := make(chan int)
    csend := make(chan bool)
    wg := new(sync.WaitGroup)
    
    go func(cdata <-chan int, csend <-chan bool, wg *sync.WaitGroup){
        for{
            select{
                case d := <-cdata:
                    total += d
                    processed++
                    wg.Done()

                case <-csend:
                    return
                    
                default:
            }
        }
    }(cdata, csend, wg)
    
    for _, d := range data{
        wg.Add(1)
        go func(d int){
            cdata <- d
        }(d)   
    }
    wg.Wait()
    csend <- true
    if total!=expectation{
        t.Fatalf("Total is %d, expected %d. Data processed %d", total, expectation, processed)
    }
}

func TestPool(t *testing.T){
    total := int(0)
    processed := int(0)
    
    alloc := alokasi.New()
    alloc.AllocationType = alokasi.AllocateAsPool
    alloc.WorkerNum = workerNum
    alloc.OnReceive = func(ac *alokasi.Context){
        ac.Allocator.Lock()
        //total = ac.Allocator.Data.Get("total", 0).(int)
        //workerid := ac.Setting.GetInt("workerid")
        intdata := ac.Data.(int)
        processed++
        total += intdata
        //toolkit.Printf("[%d w: %d] Processing data %d, Total now %d\n", processed, workerid, 
        //    intdata, total)
        //ac.Allocator.Data.Set("total", total)
        ac.Allocator.Unlock()
        //time.Sleep(100 * time.Millisecond)
    }
    alloc.Start()
    
    //wg := new(sync.WaitGroup)
    for _, d := range data{
        //wg.Add(1)
        go func(alloc *alokasi.Allocator, d int){
            //defer wg.Done()
            alloc.Send(d)
        }(alloc, d)
    }
    //wg.Wait()
    //time.Sleep(1*time.Millisecond)
    alloc.SendComplete()   
    toolkit.Println("Send all key completed")
    alloc.Wait()
    
    if total!=expectation{
        t.Fatalf("Total is %d, expected %d. Data processed %d", total, expectation, processed)
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
    
    if total!=expectation{
        t.Fatalf("Total is %d, expected %d", total, expectation)
    }
}
