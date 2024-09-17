package main

import (
	"encoding/json"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/SnowOnion/godoogle/ranking"
)

func main() {
	bytes, err := os.ReadFile("sigGraph.json")
	if err != nil {
		panic(err)
	}
	dist := make(map[ranking.SigStr]map[ranking.SigStr]int)
	err = json.Unmarshal(bytes, &dist)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(dist))
	for _, sig := range []string{
		`func(string) (int, error)`,
		`func(string) int`,
	} {
		fmt.Println()
		fmt.Println(len(dist[sig]), "~~FROM~~", sig)
		for k, v := range dist[sig] {
			fmt.Println(k, v)
		}
	}

}
