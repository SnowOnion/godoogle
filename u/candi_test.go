package u

import (
	"go/token"
	"go/types"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {
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

	///////

	cmpSliceToEle := types.NewSignatureType(nil, nil,
		[]*types.TypeParam{t1Comparable2},
		types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", types.NewSlice(t1Comparable2))),
		types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Comparable2)),
		false,
	)
	voidToVoid := types.NewSignatureType(nil, nil, nil,
		types.NewTuple(),
		types.NewTuple(),
		false)
	voidToStr := types.NewSignatureType(nil, nil, nil,
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String])),
		types.NewTuple(),
		false)
	i32IntToInt := types.NewSignatureType(nil, nil, nil,
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int32]),
			types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
		false)
	strInterfacesToVoid := types.NewSignatureType(nil, nil, nil,
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
			types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
		types.NewTuple(),
		true /*params[-1] needs to be a slice*/)
	strAnysToIntError := types.NewSignatureType(nil, nil, nil,
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
			types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
		types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
			types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())),
		true /*params[-1] needs to be a slice*/)
	anyToVoid := types.NewSignatureType(nil, nil,
		[]*types.TypeParam{t1Any},
		types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Any)),
		nil, false,
	)
	cmpToVoid := types.NewSignatureType(nil, nil,
		[]*types.TypeParam{t1Comparable},
		types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx", t1Comparable)),
		nil, false,
	)
	hsMap := types.NewSignatureType(nil, nil,
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
		false)
	loMap := types.NewSignatureType(nil, nil,
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
		false)

	// TODO more, to cover rebind*
	suite := map[string]*types.Signature{
		"[T cmp.Ordered, a,b any]func(x T, y abc.Def, z *T, ts [2]T, args ...struct{k string;  v T}) (bool, error)": nil,
		"[T cmp.Ordered]func(x T, y T) bool":                nil,
		"[S ~[]E, E constraints.Ordered] func(x S)":         nil, // https://pkg.go.dev/golang.org/x/exp/slices#Sort
		"[T comparable] func([]T)T":                         cmpSliceToEle,
		"func()":                                            voidToVoid,
		"func(string)":                                      voidToStr,
		"func(int32, int) int":                              i32IntToInt,
		"func(string,...interface{})":                       strInterfacesToVoid,
		"func (format string, a ...any) (n int, err error)": strAnysToIntError,
		"[T any] func(T)":                                   anyToVoid,
		"[T comparable] func(T)":                            cmpToVoid,
		"[a,b any] func([]a, func(a) b) []b":                hsMap,
		"[a,b any] func([]a, func(a, int) b) []b":           loMap,
	}

	i := -1 // unordered LOL
	for inp, out := range suite {
		if !slices.Contains([]string{
			"[T cmp.Ordered, a,b any]func(x T, y abc.Def, z *T, ts [2]T, args ...struct{k string;  v T}) (bool, error)",
			"[T cmp.Ordered]func(x T, y T) bool",
		}, inp) {
			continue
		}

		i++
		sig, err := Dummy1(inp)
		t.Log(i, err)
		t.Log(i, inp)
		t.Log(i, sig)
		t.Log(i, out)
		assert.Nil(t, err, "case~ %d", i)
		//assert.True(t, types.IdenticalIgnoreTags(out, sig), "case~ %d", i)
	}

}
