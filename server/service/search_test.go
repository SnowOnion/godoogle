package service

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
	"github.com/SnowOnion/godoogle/server/model"
)

// TestSearch test whether all expected answers appear in top-10 results.
func TestSearch(t *testing.T) {
	collect.InitFuncDatabase()
	ranking.DefaultRanker = ranking.NewSigGraphRanker(collect.FuncDatabase, ranking.LoadFromFile(true)) // = =、TODO be elegant!

	qna := map[string][]string{
		`func(sort.Interface)`:                               {`sort.Sort`, `sort.Stable`},
		`func(n int, cmp func(int) int) (i int, found bool)`: {`sort.Find`},
		`[T any] func(bool, T, T) T`:                         {`github.com/samber/lo.Ternary`},
		`[T any] func(bool, func() T, func() T) T`:           {`github.com/samber/lo.TernaryF`},
		`[a,b any] func([]a, func(a, int) b) []b`:            {`github.com/samber/lo.Map`},
		`[a any, b any] func([]a, func(a, int) b) []b`:       {`github.com/samber/lo.Map`}, // how do you turn this on?
		`func (string) int`:                                  {`unicode/utf8.RuneCountInString`, `strconv.Atoi`},
		`func (string) (int, error)`:                         {`strconv.Atoi`},
		`[T any] func (f func() T) <-chan T`:                 {`github.com/samber/lo.Async`, `github.com/samber/lo.Async1`},
	}
	// Collect known bad case here.
	qnaFail := map[string][]string{
		`[T any] func(collection []T, size int) [][]T`: {`github.com/samber/lo.Chunk`}, // lo.Chunk changed! func[T any, Slice ~[]T](collection Slice, size int) []Slice
		`[b,a any] func([]a, func(a, int) b) []b`:      {`github.com/samber/lo.Map`},   // consider sortTypeParamList. 函統网, UniSig.
	}

	ctx := context.Background()
	for q, as := range qna {
		resp, err := Search(ctx, model.SearchReq{Query: q})
		assert.Nil(t, err)

		top10 := resp.Result[:min(10, len(resp.Result))]
		//for _, r := range top10 {
		//	t.Log(r)
		//}
		for _, a := range as {
			assert.True(t, slices.ContainsFunc(top10, func(r model.ResultItem) bool {
				return r.FullName == a
			}))
		}
	}
	for q, as := range qnaFail {
		resp, err := Search(ctx, model.SearchReq{Query: q})
		assert.Nil(t, err)

		top10 := resp.Result[:min(10, len(resp.Result))]
		//for _, r := range top10 {
		//	t.Log(r)
		//}
		for _, a := range as {
			assert.False(t, slices.ContainsFunc(top10, func(r model.ResultItem) bool {
				return r.FullName == a
			}))
		}
	}

}
