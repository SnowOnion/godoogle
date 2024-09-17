package ranking

import (
	"fmt"
	"testing"

	"github.com/samber/lo"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/u"
)

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
		fmt.Println(ranker.Rank(lo.T2(u.Dummy(inp)).A, sigs))

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

	suite := []lo.Tuple2[string, string]{
		//`func (int,int) int`,
		//`func (int,int) bool`,
		//`[T comparable] func (T,T) bool`,
		// import...... guess what user means
		//`[S ~[]E, E constraints.Ordered] func(x S)`, // https://pkg.go.dev/golang.org/x/exp/slices#Sort
		//`[TS ~[]T, T constraints.Ordered] func(x TS)`,
		//`[E constraints.Ordered] func(x []E)`,
		//`[E constraints.Ordered] func([]E)`,
		//`[E constraints.Ordered] func([]E) []E`,
		lo.T2(`[a, b any] func (collection []a, iteratee func(item a, index int) b) []b`, `github.com/samber/lo.Map`),
		lo.T2(`[b ,a any] func (collection []a, iteratee func(item a, index int) b) []b`, `github.com/samber/lo.Map`), // TODO
		//`[T, R any] func (collection []T, iteratee func(item T, index int) R) []R`,
		//`[T, R any, K comparable] func (collection []T, iteratee func(item T, index int) R) []R`,
	}

	ranker := NaiveRanker{}

	//fmt.Println("Before ranking~ fmt.Println(sigs)~~~~", sigs)
	//for ind, sigDecl := range sigs {
	//	//sig,decl:=sigDecl.A,sigDecl.B
	//	fmt.Println(ind, sigDecl)
	//}

	for _, inOut := range suite {
		inpSig := lo.T2(u.Dummy(inOut.A)).A
		//fmt.Println(types.IdenticalIgnoreTags(q.A, q.A), q)
		fmt.Println("~~~Search result of", inpSig)
		//fmt.Println(ranker.Rank(inpSig, sigs))

		result := ranker.Rank(inpSig, sigs)
		for ind, sigDecl := range result[:min(10, len(result))] {
			fmt.Printf("%4d %s - %s\n", ind, sigDecl.B.Name(), sigDecl)
		}
		fmt.Println("~~~~~~~~~~~~")
		if inOut.B != "" {
			if result[0].B.FullName() != inOut.B {
				t.Errorf("Top 1 is %s; expected %s", result[0].B.FullName(), inOut.B)
			}
		}

		//fmt.Println("Rank~~~~", inpSig)
		//fmt.Println("Candi~~~", sigs[152].A)
		//eq := types.IdenticalIgnoreTags(sigs[152].A, inpSig)
		//fmt.Println(eq)
	}
	//fmt.Println(lo.Map(ranker.Rank(sigs[0].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[1].A, sigs), fst))
	//fmt.Println(lo.Map(ranker.Rank(sigs[2].A, sigs), fst))

}
