package ranking

import (
	"go/types"
	"sort"

	"github.com/SnowOnion/godoogle/u"
)

type NaiveRanker struct{}

// Rank by “discrete metric”(离散度量): signatures identical to query are ranked first; others are ranked arbitrarily.
func (r NaiveRanker) Rank(query *types.Signature, candidates []u.T2) []u.T2 {
	if query == nil {
		panic("Rank query == nil")
	}

	result := make([]u.T2, len(candidates))
	copy(result, candidates)
	less := func(i, j int) bool {
		candi := result[i].A
		eq := types.IdenticalIgnoreTags(candi, query)
		return eq
	}
	sort.Slice(result, less)
	return result
}
