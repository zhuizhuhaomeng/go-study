package main

import "fmt"
import "os"
import "log"
//import "runtime"
import "runtime/debug"
import   (
   "net/http"
   _ "net/http/pprof"
)

type d2array struct {
    arr [2][3] int;
}

func heapdumpHandler() {
	f, err := os.Create("metadump")
	if err != nil {
		panic(err)
	}
	//runtime.GC()
	debug.WriteHeapDump(f.Fd())
	f.Close()
}

func main() {
    var d2 d2array;
    d2.arr[1][2] = 1;
    fmt.Println(d2.arr[1][2])
    heapdumpHandler()
    go func() {
	log.Println(http.ListenAndServe("0.0.0.0:8081", nil))
    }()
    select {
    }
}
