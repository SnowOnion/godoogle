package ranking

import (
	"go/token"
	"go/types"
	"math"
	"os"
	"testing"

	//"github.com/SnowOnion/graph"
	//"github.com/SnowOnion/graph/draw"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/u"
)

func TestTrySortTypeParams(t *testing.T) {
	inps := []string{
		//"[T comparable] func([]T)T",
		//"func()",
		//"func(string)",
		//"func(int32, int) int",
		//"func(string,...interface{})",
		//"func (format string, a ...any) (n int, err error)",
		//"[T any] func(T)",
		//"[T comparable] func(T)",
		////"[a,b any] func (col []a, iter func(it a) b) (r1 []b)",         // hsMap
		////"[a,b any] func (col []a, iter func(it a) b) []b",              // hsMap
		//"[a,b any] func([]a, func(a) b) []b", // hsMap
		////"[a,b any] func(col []a, iter func(it a, idx int) b) (r1 []b)", // lo.Map
		////"[a,b any] func(col []a, iter func(it a, idx int) b) []b",      // lo.Map
		"[b,a any] func([]a, func(a, int) b) []b", // lo.Map, flip type params
		"[a,b any] func([]a, func(a, int) b) []b", // lo.Map,

	}

	for i, inp := range inps {
		sig, err := collect.Dummy(inp)
		t.Log(i, err)
		t.Log(i, inp)
		t.Log(i, sig)

		t.Log(i, sig.TypeParams())
		//t.Log(i, out)
		//assert.Nil(t, err, "case~~~~ %d", i)
		//assert.True(t, types.IdenticalIgnoreTags(out, sig), "case~~~~ %d", i)
	}
}

func TestTryDominikbraunGraph(t *testing.T) {
	/*
		g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic(), graph.Weighted())

		_ = g.AddVertex(1)
		_ = g.AddVertex(2)
		_ = g.AddVertex(3)
		_ = g.AddVertex(4)

		_ = g.AddEdge(1, 2, graph.EdgeWeight(42))
		_ = g.AddEdge(1, 3)
		_ = g.AddEdge(2, 3)
		_ = g.AddEdge(2, 4)
		_ = g.AddEdge(3, 4)

		file, _ := os.Create("./simple.gv")
		_ = draw.DOT(g, file)
		//_ = draw.DOT(g, file, draw.GraphAttribute("label", "my-graph"))
		// Then:
		// dot -Tsvg -Kneato -O simple.gv && open simple.gv.svg -a firefox

	*/

	hash := func(sig *types.Signature) string { return sig.String() }
	g2 := graph.New(hash, graph.Directed(), graph.Acyclic(), graph.Weighted())

	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())
	// v for void
	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	sei := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(e, i), false)
	vie := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(i, e), false)
	vei := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(e, i), false)
	si := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i), false)
	se := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(e), false)
	vi := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(i), false)
	ve := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(e), false)
	sv := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(), false)
	vv := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false)
	/*
		func(string) (int, error)
		func(string) (error, int)
		func() (int, error)
		func() (error, int)
		func(string) int
		func(string) error
		func() int
		func() error
		func(string)
		func()
	*/

	for _, t := range []*types.Signature{sie, sei, vie, vei, si, se, vi, ve, sv, vv} {
		_ = g2.AddVertex(t)
		//fmt.Println(t.String())
	}

	// 为啥 AddEdge 不接受 T 而是 K，要我自己调用 hash 啊，乌鱼子 // TODO PR func AddEdgeT
	// 大致是 BFS 地添加…… 但有的边是反的
	// TODO 哎，Go 可以 in-place 修改。那么参数可以既是输入又是输出。……
	_ = addEdge(g2, hash(sie), hash(sei), 1) // (PR)
	_ = addEdge(g2, hash(sie), hash(vie), 3) // (WP)
	_ = addRevE(g2, hash(sie), hash(se), 2)  // (WR)
	_ = addRevE(g2, hash(sie), hash(si), 2)  // (WR)
	_ = addEdge(g2, hash(sei), hash(sie), 1) // (PR)
	_ = addEdge(g2, hash(sei), hash(vei), 3) // (WP)
	_ = addRevE(g2, hash(sei), hash(si), 2)  // (WR)
	_ = addRevE(g2, hash(sei), hash(se), 2)  // (WR)
	_ = addEdge(g2, hash(vie), hash(vei), 1) // (PR)
	_ = addRevE(g2, hash(vie), hash(ve), 2)  // (WR)
	_ = addRevE(g2, hash(vie), hash(vi), 2)  // (WR)
	_ = addEdge(g2, hash(se), hash(ve), 3)   // (WP)
	_ = addRevE(g2, hash(se), hash(sv), 2)   // (WR)
	_ = addEdge(g2, hash(si), hash(vi), 3)   // (WP)
	_ = addRevE(g2, hash(si), hash(sv), 2)   // (WR)
	_ = addEdge(g2, hash(vei), hash(vie), 1) // (PR)
	_ = addRevE(g2, hash(vei), hash(vi), 2)  // (WR)
	_ = addRevE(g2, hash(vei), hash(ve), 2)  // (WR)
	_ = addRevE(g2, hash(vi), hash(vv), 2)   // (WR)
	_ = addRevE(g2, hash(ve), hash(vv), 2)   // (WR)
	_ = addEdge(g2, hash(sv), hash(vv), 3)   // (WP)

	file2, _ := os.Create("./siggraph.gv")
	_ = draw.DOT(g2, file2)
	t.Log("dot -Tsvg -O siggraph.gv && open siggraph.gv.svg -a firefox")

	// why not returning the sum(weight) together with the path? TODO
	path, err := graph.ShortestPath(g2, hash(vv), hash(sie))
	//path, err := graph.ShortestPath(g2, hash(sie), hash(vv))

	t.Log(err, dist(g2, hash(vv), hash(sie)), path)
	assert.Equal(t, 4, dist(g2, hash(sv), hash(sie)))
	assert.Equal(t, 0, dist(g2, hash(vv), hash(vv)))
	assert.Equal(t, math.MaxInt, dist(g2, hash(sie), hash(vv)))
	assert.Equal(t, 1, dist(g2, hash(sei), hash(sie)))
	assert.Equal(t, 1, dist(g2, hash(sie), hash(sei)))
}

func TestWeakenResults(t *testing.T) {
	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	si := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i), false)
	se := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(e), false)

	mutants := weakenResults(sie)
	assert.Equal(t, 2, len(mutants))
	for _, m := range mutants {
		t.Log(m.String())
	}
	assert.True(t, types.IdenticalIgnoreTags(se, mutants[0]))
	assert.True(t, types.IdenticalIgnoreTags(si, mutants[1]))
}

func TestWeakenParams(t *testing.T) {
	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	vie := types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(i, e), false)

	mutants := weakenParams(sie)
	assert.Equal(t, 1, len(mutants))
	for _, m := range mutants {
		t.Log(m.String())
	}
	assert.True(t, types.IdenticalIgnoreTags(vie, mutants[0]))
}

// Test only the preparation, not the ranking
func TestNewHooglyRanker(t *testing.T) {
	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)

	expectedVertices := u.SliceToSet([]SigStr{
		"func() (int, error)",
		"func() error",
		"func()",
		"func() int",
		"func(string) error",
		"func(string)",
		"func(string) int",
		"func(string) (int, error)",
	}...)
	type v2d1 struct {
		s SigStr
		t SigStr
		d int
	}
	expectedEdges := u.SliceToSet([]v2d1{
		{s: `func() int`, t: `func() (int, error)`, d: 2},
		{s: `func(string) (int, error)`, t: `func() (int, error)`, d: 3},
		{s: `func(string) error`, t: `func() error`, d: 3},
		{s: `func(string) error`, t: `func(string) (int, error)`, d: 2},
		{s: `func(string)`, t: `func()`, d: 3},
		{s: `func(string)`, t: `func(string) error`, d: 2},
		{s: `func(string)`, t: `func(string) int`, d: 2},
		{s: `func(string) int`, t: `func() int`, d: 3},
		{s: `func(string) int`, t: `func(string) (int, error)`, d: 2},
		{s: `func()`, t: `func() error`, d: 2},
		{s: `func()`, t: `func() int`, d: 2},
		{s: `func() error`, t: `func() (int, error)`, d: 2},
	}...)

	expectedPaths := u.SliceToSet([]v2d1{
		{s: `func(string) error`, t: `func(string) error`, d: 0},
		{s: `func(string)`, t: `func(string)`, d: 0},
		{s: `func(string) int`, t: `func(string) int`, d: 0},
		{s: `func(string) (int, error)`, t: `func(string) (int, error)`, d: 0},
		{s: `func() (int, error)`, t: `func() (int, error)`, d: 0},
		{s: `func() error`, t: `func() error`, d: 0},
		{s: `func()`, t: `func()`, d: 0},
		{s: `func() int`, t: `func() int`, d: 0},
		{s: `func() error`, t: `func() (int, error)`, d: 2},
		{s: `func(string) error`, t: `func(string) (int, error)`, d: 2},
		{s: `func(string)`, t: `func(string) int`, d: 2},
		{s: `func(string)`, t: `func() (int, error)`, d: 7},
		{s: `func(string)`, t: `func() error`, d: 5},
		{s: `func() int`, t: `func() (int, error)`, d: 2},
		{s: `func(string) error`, t: `func() error`, d: 3},
		{s: `func(string)`, t: `func(string) error`, d: 2},
		{s: `func(string) int`, t: `func() int`, d: 3},
		{s: `func(string)`, t: `func() int`, d: 5},
		{s: `func(string)`, t: `func()`, d: 3},
		{s: `func(string) (int, error)`, t: `func() (int, error)`, d: 3},
		{s: `func()`, t: `func() error`, d: 2},
		{s: `func()`, t: `func() int`, d: 2},
		{s: `func(string) int`, t: `func(string) (int, error)`, d: 2},
		{s: `func(string)`, t: `func(string) (int, error)`, d: 4},
		{s: `func(string) int`, t: `func() (int, error)`, d: 5},
		{s: `func(string) error`, t: `func() (int, error)`, d: 5},
		{s: `func()`, t: `func() (int, error)`, d: 4},
	}...)

	candi := []u.T2{
		{sie, types.NewFunc(token.NoPos, nil, "", sie)},
	}
	r := NewHooglyRanker(candi)
	t.Log(r.sigGraph.Order())
	t.Log(r.sigGraph.Size())
	file2, _ := os.Create("./siggraph.gv")
	_ = draw.DOT(r.sigGraph, file2) // then: dot -Tsvg -O siggraph.gv && open siggraph.gv.svg -a firefox

	V, _ := r.sigGraph.Order()
	E, _ := r.sigGraph.Size()
	assert.Equal(t, 8, V)
	assert.Equal(t, 12, E)

	vertices, _ := r.sigGraph.Vertices()
	assert.True(t, u.SliceToSet(vertices...).Equals(expectedVertices))
	//for _, v := range vertices {
	//	fmt.Printf("{s:`%s`,t:`%s`,d:%d},\n", v, v, 0)
	//}

	edges, _ := r.sigGraph.Edges()
	//for _, e := range edges {
	//	fmt.Printf("{s:`%s`,t:`%s`,d:%d},\n", e.Source, e.Target, e.Properties.Weight)
	//}
	edgeSet := u.SliceToSet(lo.Map(edges, func(e graph.Edge[SigStr], _ int) v2d1 {
		return v2d1{s: e.Source, t: e.Target, d: e.Properties.Weight}
	})...)
	assert.True(t, edgeSet.Equals(expectedEdges))

	//r.InitFloydWarshall(16)
	//for uv, d := range r.distCache {
	//	//t.Log(d, uv.A, "->", uv.B)
	//	//fmt.Printf("{s:`%s`,t:`%s`,d:%d},\n", uv.A, uv.B, d)
	//}
	pathSet := u.SliceToSet(lo.MapToSlice(r.distCache, func(uv lo.Tuple2[SigStr, SigStr], d int) v2d1 {
		return v2d1{s: uv.A, t: uv.B, d: d}
	})...)
	assert.True(t, pathSet.Equals(expectedPaths))
}

func BenchmarkDistance(b *testing.B) {
	collect.InitFuncDatabase()
	ranker := NewHooglyRanker(collect.FuncDatabase) // = =、TODO be elegant!

	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	si := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i), false)

	for i := 0; i < b.N; i++ {
		ranker.Distance(si, sie)
		//b.Log()
	}
}

func BenchmarkDistanceWithCache(b *testing.B) {
	collect.InitFuncDatabase()
	ranker := NewHooglyRanker(collect.FuncDatabase) // = =、TODO be elegant!

	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	si := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i), false)

	for i := 0; i < b.N; i++ {
		ranker.DistanceWithCache(si, sie)
		//b.Log()
	}
}

func BenchmarkDistanceWithFloydWarshall(b *testing.B) {
	collect.InitFuncDatabase()
	ranker := NewHooglyRanker(collect.FuncDatabase) // = =、TODO be elegant!

	s := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("string").Type())
	i := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("int").Type())
	e := types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type())

	sie := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i, e), false)
	si := types.NewSignatureType(nil, nil, nil, types.NewTuple(s), types.NewTuple(i), false)

	for i := 0; i < b.N; i++ {
		ranker.DistanceWithFloydWarshall(si, sie)
		//b.Log()
	}
}
