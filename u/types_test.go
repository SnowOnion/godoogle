package u

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonymize(t *testing.T) {
	////Dummy1 contains Anonymize
	//sig, err := Dummy1(`[a,b any, c comparable, d any] func(bool, a, a) (y a)`)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(sig)

	// copied from TestDummy
	t1Any3 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "a", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t2Any3 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "b", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	// lo.Map
	loMap := types.NewSignatureType(nil, nil,
		[]*types.TypeParam{t1Any3, t2Any3},
		types.NewTuple(
			types.NewVar(token.NoPos, nil, "xs",
				types.NewSlice(t1Any3),
			),
			types.NewVar(token.NoPos, nil, "f",
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
			types.NewVar(token.NoPos, nil, "ys",
				types.NewSlice(t2Any3),
			),
		),
		false)
	t.Log(loMap)

	sigA := Anonymize(loMap)
	t.Log(sigA)

	assert.Equal(t, `func[_T0, _T1 any]([]_T0, func(_T0, int) _T1) []_T1`, sigA.String())
}
