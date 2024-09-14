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

/*
1. ( ﾟ∀。) make sure in ranking.NewHooglyRanker, there is r.InitFloydWarshall(10) instead of r.InitFloydWarshallFromFile()
2. go run floydWrite.go
*/
func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	workers := max(runtime.NumCPU()-2, 1)
	fmt.Println("workers:", workers)

	collect.InitFuncDatabase()
	ranker := ranking.NewHooglyRanker(collect.FuncDatabase, false)
	ranker.InitFloydWarshall(workers)
	os.WriteFile("floyd-std_lo_graph-ttl-3-anonymize.json", ranker.MarshalDistCache(), 0666)
}
