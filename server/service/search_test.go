package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
	"github.com/SnowOnion/godoogle/server/model"
)

// TODO wip!
func TestSearch(t *testing.T) {
	collect.InitFuncDatabase()
	ranking.DefaultRanker = ranking.NewHooglyRanker(collect.FuncDatabase) // = =、TODO be elegant!

	//q := `[a,b any] func([]a, func(a, int) b) []b`
	//q := `func(sort.Interface)`
	//q := `func(n int, cmp func(int) int) (i int, found bool)`
	//q := `[T any] func(bool, T, T) T`
	//q := `[T any] func(bool, func() T, func() T) T`

	qs := []string{
		`func (string) int`,
		//`[T any] func (f func() T) <-chan T`,
		//`[T any] func(collection []T, size int) [][]T`,
		//`[a, b any] func (collection []a, iteratee func(item a, index int) b) []b`,
		//`[b, a any] func (collection []a, iteratee func(item a, index int) b) []b`, // TODO
	}
	ctx := context.Background()
	for _, q := range qs {
		resp, err := Search(ctx, model.SearchReq{Query: q})
		assert.Nil(t, err)
		t.Log(len(resp.Result))
		for _, r := range resp.Result[:5] {
			t.Log(r)
		}
	}

}
