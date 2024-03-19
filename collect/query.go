package collect

import (
	"errors"
	"go/types"
)

// Dummy rawQuery e.g.
// "func(rune) bool"
// "[T, R any, K comparable] func (collection []T, iteratee func(item T, index int) R) []R"
func Dummy(rawQuery string) (*types.Signature, error) {
	augmentedQuery := `package dummy
import (
	//"golang.org/x/exp/constraints"
	"sort"
	"time"
)
type dummy ` + rawQuery
	sigs, err := ParseGenDeclTypeSpecFuncSigs(augmentedQuery)
	if err != nil {
		return nil, err
	}
	if len(sigs) == 0 {
		return nil, errors.New("no type signature in augmentedQuery")
	}
	return sigs[0], nil
}

// DummyF rawQuery e.g.
// "func f(rune) bool"
// " func fff[T, R any, K comparable](collection []T, f func(item T, index int) R) []R"
func DummyF(rawQuery string) (*types.Signature, error) {
	augmentedQuery := `package dummy
	` + rawQuery + `{return}`
	sigs, err := ParseFuncSigs(augmentedQuery)

	if err != nil {
		return nil, err
	}
	if len(sigs) == 0 {
		return nil, errors.New("no type signature in augmentedQuery")
	}
	return sigs[0].A, nil
}

type exampleFunctionTypeNonGeneric func(x int) string

// > must repeat the function signature
// https://stackoverflow.com/a/9596177/2801663
var exampleFNG1 = func(x int) string {
	return ""
}
var exampleFNG2 exampleFunctionTypeNonGeneric = exampleFNG1

// Few people would write ↓ (although valid),
type exampleFunctionTypeGeneric[T, R any, K comparable] func(collection []T, iteratee func(item T, index int) R) []R

// because:
// > A function literal cannot be generic because the function literal produces a function value, and the function value cannot be generic.
// > Generic functions are not a type in the Go type system; they are just a syntactic tool for instantiating functions with type substitution.
// https://stackoverflow.com/a/75915437/2801663
// https://go.dev/ref/spec#FunctionLit
//var exampleFG1 = [T, R any, K comparable] func(collection []T, iteratee func(item T, index int) R) []R{
//
//}

var exampleFG2 = func(collection []int, iteratee func(item int, index int) string) []string {
	return nil
}
var exampleFG3 exampleFunctionTypeGeneric[int, string, bool] = exampleFG2
