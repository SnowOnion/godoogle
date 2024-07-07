package ranking

import (
	"encoding/json"
	"go/types"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/samber/lo"

	"github.com/SnowOnion/godoogle/u"
)

type HooglyRanker struct {
	sigIndex  map[SigStr][]u.T2
	hash      func(sig *types.Signature) string
	sigGraph  graph.Graph[SigStr, *types.Signature]
	distCache map[lo.Tuple2[SigStr, SigStr]]int // TODO just use map[SigStr]map[SigStr]int?
}

type SigStr = string // types.Signature#String()

func NewHooglyRanker(candidates []u.T2) HooglyRanker {
	hash := func(sig *types.Signature) string {
		return sig.String()
	}
	r := HooglyRanker{
		sigIndex:  make(map[SigStr][]u.T2),
		hash:      hash,
		sigGraph:  graph.New(hash, graph.Directed(), graph.Acyclic(), graph.Weighted()),
		distCache: make(map[lo.Tuple2[SigStr, SigStr]]int),
	}
	r.InitCandidates(candidates)

	return r
}

func (r *HooglyRanker) InitCandidates(candidates []u.T2) {
	for _, t := range candidates {
		//hlog.Debug(i, " ", t.String())
		sigStr := t.A.String() // todo anonymize?
		_, ok := r.sigIndex[sigStr]
		if ok {
			r.sigIndex[sigStr] = append(r.sigIndex[sigStr], t)
			continue
		}
		r.sigIndex[sigStr] = []u.T2{t}

		// TODO 暴露配置项
		r.InitDFS(Anonymize(t.A), 3)
	}
	hlog.Info("|candidates|=", len(candidates), "; |sigIndex|=", len(r.sigIndex))
	o, _ := r.sigGraph.Order()
	s, _ := r.sigGraph.Size()
	hlog.Info("Graph order and size: |V|=", o, "; |E|=", s)
	file2, _ := os.Create("./siggraph.gv")
	_ = draw.DOT(r.sigGraph, file2) // then: dot -Tsvg -O siggraph.gv && open siggraph.gv.svg -a firefox

	//r.InitFloydWarshallFromFile()
	//hlog.Info("Before InitFloydWarshallFromFile")
	//r.InitFloydWarshall(10)
	//hlog.Info("After InitFloydWarshallFromFile")
}

func (r *HooglyRanker) InitDFS(sig *types.Signature, depthTTL int) {
	if _, err := r.sigGraph.Vertex(r.hash(sig)); err == nil {
		return
	}
	sig = Anonymize(sig) // if not Anonymize: things like `func(s string) string` go in... // TODO 提效：一棵树只需 Anonymize 一次
	_ = r.sigGraph.AddVertex(sig)

	if depthTTL <= 0 {
		return
	}

	// pause, sleep
	//for _, mut := range permuteParams(sig) {
	//	r.InitDFS(mut)
	//	_ = addEdge(r.sigGraph, r.hash(mut), r.hash(sig), 1)
	//}
	//for _, mut := range permuteResults(sig) {
	//	r.InitDFS(mut)
	//	_ = addEdge(r.sigGraph, r.hash(mut), r.hash(sig), 1)
	//}
	for _, mut := range weakenParams(sig) {
		r.InitDFS(mut, depthTTL-1)
		_ = addEdge(r.sigGraph, r.hash(sig), r.hash(mut), 3)
	}
	for _, mut := range weakenResults(sig) {
		r.InitDFS(mut, depthTTL-1)
		_ = addEdge(r.sigGraph, r.hash(mut), r.hash(sig), 2)
	}
}

// InitFloydWarshall refresh distCache by applying Floyd-Warshall algorithm to sigGraph.
func (r *HooglyRanker) InitFloydWarshall(numWorkers int) {
	r.distCache = make(map[lo.Tuple2[SigStr, SigStr]]int)
	// It would be stupid to lock the whole distCache... Wanna use map[SigStr]map[SigStr]int now. TODO

	g := r.sigGraph
	adj, err := g.AdjacencyMap()
	if err != nil {
		panic("InitFloydWarshall .AdjacencyMap(): " + err.Error())
	}

	vertices, err := g.Vertices()
	if err != nil {
		panic("InitFloydWarshall .Vertices(): " + err.Error())
	}
	vertexIndex := make(map[SigStr]int)
	for i, v := range vertices {
		vertexIndex[v] = i
	}

	const inf = math.MaxInt
	// 空间大点大点吧，访问快 + 本算法里不用锁
	dist := make([][]int, len(vertices))
	for i := 0; i < len(dist); i++ {
		dist[i] = make([]int, len(vertices))
		for j := 0; j < len(dist[i]); j++ {
			if i != j {
				dist[i][j] = inf
			}
		}
	}

	for u, es := range adj {
		i := vertexIndex[u]
		for v, e := range es {
			j := vertexIndex[v]
			dist[i][j] = e.Properties.Weight
		}
	}

	var wg sync.WaitGroup
	chunks := len(vertices) / numWorkers
	for ki, k := range vertices {
		hlog.Infof("%d/%d k=%s", ki+1, len(vertices), k)
		wg.Add(numWorkers)
		for worker := 0; worker < numWorkers; worker++ {
			go func(worker int) {
				defer wg.Done()
				start := worker * chunks
				end := lo.Ternary(worker == numWorkers-1, len(vertices), start+chunks)

				for i := start; i < end; i++ {
					for j := range vertices {
						if dIK, dKJ := dist[i][ki], dist[ki][j]; dIK != inf && dKJ != inf {
							if dist[i][j] > dIK+dKJ {
								dist[i][j] = dIK + dKJ
							}
						}

					}
				}
			}(worker)
		}
		wg.Wait()
	}
	for i := 0; i < len(dist); i++ {
		for j := 0; j < len(dist[0]); j++ {
			if dist[i][j] != inf {
				r.distCache[lo.T2(vertices[i], vertices[j])] = dist[i][j]
			}
		}
	}
}

func (r *HooglyRanker) InitFloydWarshallFromFile() {
	path := "res/floyd.json" // TODO elegant
	if len(os.Args) >= 2 && os.Args[1] != "" {
		path = os.Args[1]
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dist := make(map[SigStr]map[SigStr]int)
	err = json.Unmarshal(bytes, &dist)
	if err != nil {
		panic(err)
	}

	for u, es := range dist {
		for v, d := range es {
			r.distCache[lo.T2(u, v)] = d
		}
	}
}

func (r *HooglyRanker) MarshalDistCache() []byte {
	// map map to map
	m := make(map[SigStr]map[SigStr]int)
	for uv, d := range r.distCache {
		if _, ok := m[uv.A]; !ok {
			m[uv.A] = make(map[SigStr]int)
		}
		m[uv.A][uv.B] = d
	}

	j, err := json.Marshal(m)
	if err != nil {
		panic("MarshalDistCache err=" + err.Error())
	}
	return j

}

func (r *HooglyRanker) UnmarshalDistCache(j []byte) {

}

// Anonymize 先写一层吧……累了。显而易见的 badcase: lo.Map 模糊搜索不到 TODO recursive
func Anonymize(sig *types.Signature) *types.Signature {
	return types.NewSignatureType(
		anonymizeVar(sig.Recv()),
		copyTparams(typeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
		copyTparams(typeParamListToSliceOfTypeParam(sig.TypeParams())),
		types.NewTuple(lo.Map(tupleToSliceOfVar(sig.Params()), anonymizeVarI)...),
		types.NewTuple(lo.Map(tupleToSliceOfVar(sig.Results()), anonymizeVarI)...),
		sig.Variadic(),
	)
}
func anonymizeVarI(v *types.Var, _ int) *types.Var {
	return anonymizeVar(v)
}
func anonymizeVar(v *types.Var) *types.Var {
	if v == nil {
		return v
	}
	return types.NewVar(v.Pos(), nil, "", v.Type())
}

func (r *HooglyRanker) Distance(src, tar *types.Signature) int {
	return dist(r.sigGraph, r.hash(src), r.hash(tar))
}

func (r *HooglyRanker) DistanceWithCache(src, tar *types.Signature) int {
	key := lo.T2(r.hash(src), r.hash(tar))
	if d, ok := r.distCache[key]; ok {
		return d
	}
	r.distCache[key] = dist(r.sigGraph, r.hash(src), r.hash(tar))
	return r.distCache[key]
}

func (r *HooglyRanker) DistanceWithFloydWarshall(src, tar *types.Signature) int {
	key := lo.T2(r.hash(src), r.hash(tar))
	if d, ok := r.distCache[key]; ok {
		return d
	}
	return math.MaxInt
}

//func permuteParams(sig *types.Signature) []*types.Signature {
//	return nil
//}

// (x1,x2,x3) -> Y  ~~>
// [(x2,x3) -> Y, (x1,x3) -> Y, (x2,x3) -> Y]
func weakenParams(sig *types.Signature) []*types.Signature {
	rst := make([]*types.Signature, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		newParamsSlice := tupleToSliceOfVarExcept(sig.Params(), i)
		newParams := types.NewTuple(newParamsSlice...)

		newVariadic := sig.Variadic()
		if i == sig.Params().Len()-1 {
			newVariadic = false // otherwise: panic: got int, want variadic parameter with unnamed slice type or string as core type
		}

		rst[i] = types.NewSignatureType(sig.Recv(),
			copyTparams(typeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
			copyTparams(typeParamListToSliceOfTypeParam(sig.TypeParams())), // if not copy -> panic: type parameter bound more than once
			newParams,
			sig.Results(),
			newVariadic)
	}
	return rst
}

// X -> (y1,y2,y3)  ~~>
// [X -> (y2,y3), X -> (y1,y3), X -> (y2,y3)]
func weakenResults(sig *types.Signature) []*types.Signature {
	rst := make([]*types.Signature, sig.Results().Len())
	for i := 0; i < sig.Results().Len(); i++ {
		newResultSlice := tupleToSliceOfVarExcept(sig.Results(), i)
		newResult := types.NewTuple(newResultSlice...)
		rst[i] = types.NewSignatureType(sig.Recv(),
			copyTparams(typeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
			copyTparams(typeParamListToSliceOfTypeParam(sig.TypeParams())), // if not copy -> panic: type parameter bound more than once
			sig.Params(),
			newResult,
			sig.Variadic())
	}
	return rst
}

func addEdge[K comparable, T any](g graph.Graph[K, T], src, tar K, weight int) error {
	return g.AddEdge(src, tar, graph.EdgeWeight(weight), graph.EdgeAttribute("label", strconv.Itoa(weight)))
}

func addRevE[K comparable, T any](g graph.Graph[K, T], tar, src K, weight int) error {
	return addEdge(g, src, tar, weight)
}

// if src==tar, 0;
// if not reachable, math.MaxInt;
// else, sum(weight) of the shortest path.
func dist[K comparable, T any](g graph.Graph[K, T], src, tar K) int {
	d := math.MaxInt
	path, err := graph.ShortestPath(g, src, tar)

	if err == nil {
		d = 0
		for i := range path[:len(path)-1] {
			edge, err2 := g.Edge(path[i], path[i+1])
			if err2 != nil {
				// should not happen TODO
				continue
			} else {
				d += edge.Properties.Weight
			}
		}
	}
	return d
}

// Rank by distance
// TODO remove candidates param
func (r HooglyRanker) Rank(query *types.Signature, candidates []u.T2) []u.T2 {
	if query == nil {
		panic("Rank query == nil")
	}
	query = Anonymize(query)

	result := make([]u.T2, len(candidates))
	copy(result, candidates)
	for i := 0; i < len(result); i++ {
		result[i].A = Anonymize(result[i].A)
	}
	//less := func(i, j int) bool {
	//	return distance(query, result[i].A) < distance(query, result[j].A)
	//}

	//// debug
	//for i, candidate := range result {
	//	hlog.Debugf("Distance!%d %d %s", i, r.Distance(query, candidate.A), candidate.A)
	//}

	less := func(i, j int) bool {
		return r.DistanceWithFloydWarshall(query, result[i].A) < r.DistanceWithFloydWarshall(query, result[j].A)
	}
	sort.Slice(result, less)
	return result
}

// distance src -> dst
// distance(s,r) may != distance(r,s)
// TODO memoize
// TODO buggy: bind more than once.
func distance(src, dst *types.Signature) int64 {
	if eq := types.IdenticalIgnoreTags(src, dst); eq {
		return 0
	}

	//permute type params
	srcSortTparams := signatureWithNewTypeParamList(src, sortTypeParamList(src.TypeParams()))
	dstSortTparams := signatureWithNewTypeParamList(dst, sortTypeParamList(dst.TypeParams()))
	if eq := types.IdenticalIgnoreTags(srcSortTparams, dstSortTparams); eq {
		return 1
	}

	// permute params
	//src.Params()

	// others.
	return 114514
}

func typeParamListToSliceOfTypeParam(inp *types.TypeParamList) []*types.TypeParam {
	out := make([]*types.TypeParam, inp.Len())
	for i := 0; i < inp.Len(); i++ {
		out[i] = inp.At(i)
	}
	return out
}

func sortTypeParamList(inp *types.TypeParamList) []*types.TypeParam {
	out := typeParamListToSliceOfTypeParam(inp)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Obj().Name() < out[j].Obj().Name()
	})
	return out
}

// very immutable! very purely functional
func signatureWithNewTypeParamList(inp *types.Signature, tparams []*types.TypeParam) *types.Signature {
	return types.NewSignatureType(inp.Recv(),
		copyTparams(typeParamListToSliceOfTypeParam(inp.RecvTypeParams())), // TODO no copy but no panic? Oh I have not bound (rebind) them for the first time.
		copyTparams(tparams), // if not copy -> panic: type parameter bound more than once
		inp.Params(),
		inp.Results(),
		inp.Variadic())
}

func copyTparams(tparams []*types.TypeParam) []*types.TypeParam {
	cp := make([]*types.TypeParam, len(tparams))
	for i := 0; i < len(tparams); i++ {
		tp := tparams[i]
		cp[i] = types.NewTypeParam(tp.Obj(), tp.Constraint())
	}
	return cp
}

func tupleToSliceOfVar(inp *types.Tuple) []*types.Var {
	out := make([]*types.Var, inp.Len())
	for i := 0; i < inp.Len(); i++ {
		out[i] = inp.At(i)
	}
	return out
}

// todo see also https://pkg.go.dev/golang.org/x/exp/slices#Delete
func tupleToSliceOfVarExcept(inp *types.Tuple, except int) []*types.Var {
	out := make([]*types.Var, 0, inp.Len()-1)
	for i := 0; i < inp.Len(); i++ {
		if i != except {
			out = append(out, inp.At(i))
		}
	}
	return out
}
