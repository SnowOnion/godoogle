package ranking

import (
	"go/types"

	"github.com/SnowOnion/godoogle/u"
)

type Ranker interface {
	// Rank ranks candidates by their relevance to query, like a search engine.
	//
	// Rank does NOT mutate [candidates].
	Rank(query *types.Signature, candidates []u.T2) []u.T2
}
type IdentityRanker struct{}

// Rank simply returns the reference of [candidates].
func (r IdentityRanker) Rank(query *types.Signature, candidates []u.T2) []u.T2 {
	return candidates
}

var DefaultRanker Ranker = HooglyRanker{}
