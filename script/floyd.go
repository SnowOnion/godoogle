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
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	collect.InitFuncDatabase()
	ranker := ranking.NewHooglyRanker(collect.FuncDatabase)
	ranker.InitFloydWarshall()
	os.WriteFile("floyd-std-ttl-3.json", ranker.MarshalDistCache(), 0666)
}
