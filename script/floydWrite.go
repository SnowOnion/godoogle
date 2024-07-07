package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	workers := max(runtime.NumCPU()-2, 1)
	fmt.Println("workers:", workers)

	collect.InitFuncDatabase()
	ranker := ranking.NewHooglyRanker(collect.FuncDatabase)
	ranker.InitFloydWarshall(workers)
	os.WriteFile("floyd-std_lo_graph-ttl-3.json", ranker.MarshalDistCache(), 0666)
}
