package u

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {
	// TODO more, to cover rebind*
	inps := []string{
		"[T comparable] func([]T)T",
		"func()",
		"func(string)",
		"func(int32, int) int",
		"func(string,...interface{})",
		"func (format string, a ...any) (n int, err error)",
		"[T any] func(T)",
		"[T comparable] func(T)",
		//"[a,b any] func (col []a, iter func(it a) b) (r1 []b)",         // hsMap
		//"[a,b any] func (col []a, iter func(it a) b) []b",              // hsMap
		"[a,b any] func([]a, func(a) b) []b", // hsMap
		//"[a,b any] func(col []a, iter func(it a, idx int) b) (r1 []b)", // lo.Map
		//"[a,b any] func(col []a, iter func(it a, idx int) b) []b",      // lo.Map
		"[a,b any] func([]a, func(a, int) b) []b",   // lo.Map,
		"[S ~[]E, E constraints.Ordered] func(x S)", // https://pkg.go.dev/golang.org/x/exp/slices#Sort
		"[E constraints.Ordered] func(x []E)",
	}

	t1Any := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t1Comparable := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil),
		types.Universe.Lookup("comparable").Type())
	t1Comparable2 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil),
		types.Universe.Lookup("comparable").Type())
	t1Any2 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t2Any2 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t1Any3 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t2Any3 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())

	outputs := []*types.Signature{
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Comparable2},
			types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", types.NewSlice(t1Comparable2))),
			types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Comparable2)),
			false,
		),
		types.NewSignatureType(nil, nil, nil,
			types.NewTuple(),
			types.NewTuple(),
			false),
		types.NewSignatureType(nil, nil, nil,
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String])),
			types.NewTuple(),
			false),
		types.NewSignatureType(nil, nil, nil,
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int32]),
				types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
			false),
		types.NewSignatureType(nil, nil, nil,
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
				types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
			types.NewTuple(),
			true /*params[-1] needs to be a slice*/),
		types.NewSignatureType(nil, nil, nil,
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
				types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
			types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
				types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())),
			true /*params[-1] needs to be a slice*/),
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Any},
			types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Any)),
			nil, false,
		),
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Comparable},
			types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Comparable)),
			nil, false,
		),

		// hsMap
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Any2, t2Any2},
			types.NewTuple(
				types.NewVar(token.NoPos, nil, "",
					types.NewSlice(t1Any2),
				),
				types.NewVar(token.NoPos, nil, "",
					types.NewSignatureType(nil, nil, nil,
						types.NewTuple(
							types.NewVar(token.NoPos, nil, "", t1Any2),
						),
						types.NewTuple(types.NewVar(token.NoPos, nil, "", t2Any2)),
						false),
				),
			),
			types.NewTuple(
				types.NewVar(token.NoPos, nil, "",
					types.NewSlice(t2Any2),
				),
			),
			false),

		// lo.Map
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Any3, t2Any3},
			types.NewTuple(
				types.NewVar(token.NoPos, nil, "",
					types.NewSlice(t1Any3),
				),
				types.NewVar(token.NoPos, nil, "",
					types.NewSignatureType(nil, nil, nil,
						types.NewTuple(
							types.NewVar(token.NoPos, nil, "", t1Any3),
							types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]), // lo > hs
						),
						types.NewTuple(types.NewVar(token.NoPos, nil, "", t2Any3)),
						false),
				),
			),
			types.NewTuple(
				types.NewVar(token.NoPos, nil, "",
					types.NewSlice(t2Any3),
				),
			),
			false),
		nil,
		nil,
	}

	for i, inp := range inps {
		out := outputs[i]
		sig, err := Dummy(inp)
		t.Log(i, err)
		t.Log(i, inp)
		t.Log(i, sig)
		t.Log(i, out)
		assert.Nil(t, err, "case~~~~ %d", i)
		assert.True(t, types.IdenticalIgnoreTags(out, sig), "case~~~~ %d", i)
	}

}
