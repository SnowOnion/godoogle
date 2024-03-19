package ranking

import (
	"github.com/SnowOnion/godoogle/u"
	"go/types"
	"sort"
)

type NaiveRanker struct{}

/*
Heuristics
(Mixing Go and Haskell notations...)
Distances:
0. a -> b ~~ a -> b // identical
+1. a -> b -> c ~~ b -> a -> c // flip params
	+1. (a,b,c) -> d ~~ (c,a,b) -> d // flip params twice
+1. a -> (b,c) ~~ a -> (c,b) // flip results
	+1. a -> (b,c,d) ~~ a -> (c,d,b) // flip results twice

+2. a -> b ~~ a -> c // possible first step
+2. a -> b ~~ c -> b // possible last step
---
- subtyping


*/

func (r NaiveRanker) Rank(query *types.Signature, candidates []u.T2) []u.T2 {
	if query == nil {
		panic("Rank query == nil")
	}

	result := make([]u.T2, len(candidates))
	copy(result, candidates)
	less := func(i, j int) bool {
		candi := result[i].A
		eq := types.IdenticalIgnoreTags(candi, query)
		//fmt.Printf("%s==%s: %t\n", candi, query, eq)
		return eq
	}
	sort.Slice(result, less)
	return result
}
