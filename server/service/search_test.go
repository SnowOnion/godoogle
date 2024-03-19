package service

import (
	"context"
	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/server/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	collect.InitFuncDatabase()

	//q := `[a,b any] func([]a, func(a, int) b) []b`
	//q := `func(sort.Interface)`
	//q := `func(n int, cmp func(int) int) (i int, found bool)`
	//q := `[T any] func(bool, T, T) T`
	//q := `[T any] func(bool, func() T, func() T) T`

	qs := []string{
		`[T any] func (f func() T) <-chan T`,
		`[T any] func(collection []T, size int) [][]T`,
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
