package collect

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/token"
	"go/types"
	"testing"
)

func Test1(t *testing.T) {

	src := `
package main

import "fmt"

// Add adds two integers and returns the result.
func Add(a, b int) int {
	return a + b
}

func Eq[T comparable](a, b T) bool {
	return a == b
}

// PrintHello prints a hello message.
func PrintHello(name string) {
	fmt.Println("Hello,", name)
}
`
	sigs, err := ParseFuncSigs(src)
	if err != nil {
		t.Error(err)
	}

	for _, sig := range sigs {
		fmt.Println(sig)
	}

}

func Test2(t *testing.T) {
	//path := `/Users/snowonion/develop/Golang/src/github.com/samber/lo`
	pkgID := `github.com/samber/lo`
	sigs, err := ParseFuncSigsFromPackage(pkgID)
	if err != nil {
		t.Error(err)
	}

	for _, sig := range sigs {
		fmt.Println(sig)
	}
}

func Test3(t *testing.T) {
	//path := `/Users/snowonion/develop/Golang/src/github.com/samber/lo`
	path := `github.com/samber/lo`
	LoadDirDoc(path)
	//sigs, err :=
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//for _, sig := range sigs {
	//	fmt.Println(sig)
	//}
}

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

// !         	Error:      	Expected nil, but got: types.Error{Fset:(*token.FileSet)(0xc00016b280), Pos:41, Msg:"missing return", Soft:false, go116code:102, go116start:41, go116end:41}
// 补上 named return 和 {return} 呢…… -> 可以！但要注意定义和引用相同的 TypeParam 对象~~~
func TestDummyF(t *testing.T) {
	inps := []string{
		//"func f()",
		//"func ff(string)",
		//"func fff(int32, int) (r1 int)",
		//"func ffff (string,...interface{})",
		//"func fffff (format string, a ...any) (n int, err error)",
		"func ft[T any](T)",
		"func lomapx[a comparable, b any](collection []a, iteratee func(item a, index int) b) (r1 []b)",
	}
	inps2 := []string{
		"func f()",
		"func ff(string)",
		"func fff(int32, int) int",
		"func ffff (string,...interface{})",
		"func fffff (format string, a ...any) (n int, err error)",
	}
	inps2 = inps2

	t1Any := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t1AnyToo := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	t1AnyToo = t1AnyToo
	t1Comparable2 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T1", nil /*这里？TODO*/),
		types.Universe.Lookup("comparable").Type())
	t2Any2 := types.NewTypeParam(
		types.NewTypeName(token.NoPos, nil, "T2", nil /*这里？TODO*/),
		types.Universe.Lookup("any").Type())
	// 可以（且必须？）用于同一个 signature，但不能复用于多个 signature! 否则在 sig.tparams = bindTParams(typeParams)
	// 的时候， typ.index >= 0  ->	panic: type parameter bound more than once

	outputs := []*types.Signature{
		//types.NewSignatureType(nil, nil, nil,
		//	types.NewTuple(),
		//	types.NewTuple(),
		//	false),
		//types.NewSignatureType(nil, nil, nil,
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String])),
		//	types.NewTuple(),
		//	false),
		//types.NewSignatureType(nil, nil, nil,
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int32]),
		//		types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int])),
		//	false),
		//types.NewSignatureType(nil, nil, nil,
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
		//		types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
		//	types.NewTuple(),
		//	true /*params[-1] needs to be a slice*/),
		//types.NewSignatureType(nil, nil, nil,
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String]),
		//		types.NewVar(token.NoPos, nil, "", types.NewSlice(types.NewInterfaceType(nil, nil)) /*any/interface{}*/)),
		//	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
		//		types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())),
		//	true /*params[-1] needs to be a slice*/),
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Any},
			types.NewTuple(types.NewVar(token.NoPos, nil, "xxxx" /*新构造一个一样的，就不 Identical 了！*/, t1Any)),
			nil, false,
		),
		types.NewSignatureType(nil, nil,
			[]*types.TypeParam{t1Comparable2, t2Any2},
			types.NewTuple(
				types.NewVar(token.NoPos, nil, "",
					types.NewSlice(t1Comparable2),
				),
				types.NewVar(token.NoPos, nil, "",
					types.NewSignatureType(nil, nil, nil,
						types.NewTuple(
							types.NewVar(token.NoPos, nil, "", t1Comparable2),
							types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
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
	}

	for i, inp := range inps {
		sig, err := DummyF(inp)
		t.Log(i, err)
		t.Log(i, inp)
		t.Log(i, sig)
		out := outputs[i]
		t.Log(i, out)
		assert.Nil(t, err, "case~~~~ %d", i)
		assert.True(t, types.IdenticalIgnoreTags(out, sig), "case~~~~ %d", i)
	}

	/*
		3 parties:
		1. from query->dummy
			type style
			func style
		2. from types.New*
		3. from database
			Go package source
			other
	*/
}

type A[T any] struct{}

func (A[int]) m() {}

//func f[T any](func(func(T)), A[T], T[T]) {} // invalid operation: T[T] (T is not a generic type)

//var g = f[int]
