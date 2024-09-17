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
	// pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	workers := max(runtime.NumCPU()-2, 1)
	fmt.Println("workers:", workers)

	collect.InitFuncDatabase()
	ranker := ranking.NewSigGraphRanker(collect.FuncDatabase, ranking.LoadFromFile(false))
	ranker.InitFloydWarshall(workers)

	if err := os.WriteFile("sigGraph.json", ranker.MarshalDistCache(), 0666); err != nil {
		panic(err)
	}
}
