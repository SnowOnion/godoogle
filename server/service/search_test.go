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
	ranking.DefaultRanker = ranking.NewHooglyRanker(collect.FuncDatabase, ranking.LoadFromFile(true)) // = =、TODO be elegant!

	qna := map[string][]string{
		`func(sort.Interface)`:                               {`sort.Sort`, `sort.Stable`},
		`func(n int, cmp func(int) int) (i int, found bool)`: {`sort.Find`},
		`[T any] func(bool, T, T) T`:                         {`github.com/samber/lo.Ternary`},
		`[T any] func(bool, func() T, func() T) T`:           {`github.com/samber/lo.TernaryF`},
		`[a,b any] func([]a, func(a, int) b) []b`:            {`github.com/samber/lo.Map`},
		`func (string) int`:                                  {`unicode/utf8.RuneCountInString`, `strconv.Atoi`},
		`func (string) (int, error)`:                         {`strconv.Atoi`},
		`[T any] func (f func() T) <-chan T`:                 {`github.com/samber/lo.Async`, `github.com/samber/lo.Async1`},
		//`[T any] func(collection []T, size int) [][]T`:       {``}, // lo.Chunk changed!
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

}
