package u

import (
	"go/types"
	"maps"

	"github.com/samber/lo"
)

type T2 lo.Tuple2[*types.Signature, *types.Func]

//type T2T lo.Tuple2[*types.Signature, *ast.GenDecl]

// for better debugging
func (t T2) String() string {
	return t.B.FullName() + " :: " + t.A.String()
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

type Set[T comparable] map[T]struct{}

func SliceToSet[T comparable](xs ...T) Set[T] {
	s := make(Set[T])
	e := struct{}{}
	for _, x := range xs {
		s[x] = e
	}
	return s
}

func (s Set[T]) Contains(x T) bool {
	_, ok := s[x]
	return ok
}

func (s Set[T]) Equals(t Set[T]) bool {
	return maps.Equal(s, t)
}
