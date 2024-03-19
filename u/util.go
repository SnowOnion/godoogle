package u

import (
	"github.com/samber/lo"
	"go/types"
)

type T2 lo.Tuple2[*types.Signature, *types.Func]

//type T2T lo.Tuple2[*types.Signature, *ast.GenDecl]

// for better debugging
func (t T2) String() string {
	return t.B.FullName() + ": " + t.A.String()
}

// KsVs returns []K, []V in corresponding order.
func KsVs[K comparable, V any](m map[K]V) ([]K, []V) {
	ks := make([]K, len(m))
	vs := make([]V, len(m))
	i := 0
	for k, v := range m {
		ks[i] = k
		vs[i] = v
	}
	return ks, vs
}
