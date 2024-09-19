package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
)

func main() {
	// pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	collect.InitFuncDatabase()
	ranker := ranking.NewSigGraphRanker(collect.FuncDatabase, ranking.LoadFromFile(false))

	if err := os.WriteFile("sigGraph.json", ranker.MarshalDistCache(), 0666); err != nil {
		panic(err)
	}
}
