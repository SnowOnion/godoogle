package ranking

import (
	"fmt"
	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/u"
	"github.com/samber/lo"
	"go/types"
	"testing"
)

type a func(int, int) int
type b func(int, int) bool
type c[T comparable] func(T, T) bool

func TestNaive1(t *testing.T) {
	src := `
package main

import "fmt"

// PrintHello prints a hello message.
func PrintHello(name string) {
	fmt.Println("Hello,", name)
}

// Add adds two integers and returns the result.
func Add(a, b int) int {
	return a + b
}

func Eq[T comparable](a, b T) bool {
	return a == b
}


`
	inps := []string{
		`func (int,int) int`,
		`func (int,int) bool`,
		`[T comparable] func (T,T) bool`,
	}

	sigs, err := collect.ParseFuncSigs(src)
	if err != nil {
		t.Error(err)
	}

	ranker := NaiveRanker{}

	fmt.Println(sigs)
	for _, inp := range inps {
		//fmt.Println(types.IdenticalIgnoreTags(q.A, q.A), q)
		fmt.Println()
		fmt.Println(ranker.Rank(lo.T2(collect.Dummy(inp)).A, sigs))

	}
	//fmt.Println(lo.Map(ranker.Rank(sigs[0].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[1].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[2].A, sigs), fst))

}

func TestNaive2(t *testing.T) {
	pkgIDs := []string{
		`github.com/samber/lo`,
		`sort`,
		`golang.org/x/exp/slices`,
	}
	sigs, err := collect.ParseFuncSigsFromPackage(pkgIDs...)
	if err != nil {
		t.Error(err)
	}

	inps := []string{
		//`func (int,int) int`,
		//`func (int,int) bool`,
		//`[T comparable] func (T,T) bool`,
		// import...... guess what user means
		//`[S ~[]E, E constraints.Ordered] func(x S)`, // https://pkg.go.dev/golang.org/x/exp/slices#Sort
		//`[TS ~[]T, T constraints.Ordered] func(x TS)`,
		//`[E constraints.Ordered] func(x []E)`,
		//`[E constraints.Ordered] func([]E)`,
		//`[E constraints.Ordered] func([]E) []E`,
		`[a, b any] func (collection []a, iteratee func(item a, index int) b) []b`, // lo.Map
		//`[T, R any] func (collection []T, iteratee func(item T, index int) R) []R`,
		//`[T, R any, K comparable] func (collection []T, iteratee func(item T, index int) R) []R`,
	}

	ranker := NaiveRanker{}

	//fmt.Println("Before ranking~ fmt.Println(sigs)~~~~", sigs)
	//for ind, sigDecl := range sigs {
	//	//sig,decl:=sigDecl.A,sigDecl.B
	//	fmt.Println(ind, sigDecl)
	//}

	for _, inp := range inps {
		inpSig := lo.T2(collect.Dummy(inp)).A
		//fmt.Println(types.IdenticalIgnoreTags(q.A, q.A), q)
		fmt.Println("~~~Search result of", inpSig)
		//fmt.Println(ranker.Rank(inpSig, sigs))

		for ind, sigDecl := range ranker.Rank(inpSig, sigs) {
			fmt.Printf("%4d %s#%s\n", ind, sigDecl.B.Name(), sigDecl)
		}
		fmt.Println("~~~~~~~~~~~~")

		//fmt.Println("Rank~~~~", inpSig)
		//fmt.Println("Candi~~~", sigs[152].A)
		//eq := types.IdenticalIgnoreTags(sigs[152].A, inpSig)
		//fmt.Println(eq)
	}
	//fmt.Println(lo.Map(ranker.Rank(sigs[0].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[1].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[2].A, sigs), fst))

}

// func fst[A, B any](t lo.Tuple2[A, B], _ int) A { return t.A }
func fst(t u.T2, _ int) *types.Signature { return t.A }
