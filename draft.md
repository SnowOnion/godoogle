
```
# Heuristics
(Mixing Go and Haskell notations...)

从类型 f 指向类型 g，则用 g 的实例可以实现 f 的实例。（考虑 柯里-霍华德对应。）
例：
// func(string) (int, error) -> func(string) int
// 一种实现方式：
func Atoi(s string) (int, error);
func MustAtoi(s string) int {
    i, _ := Atoi(s)
    return i
}

## Distances:
0.
	a -> b <~~> a -> b // identical (ID)
1.
	(x1,x2) -> (y1,y2) <~~> (x2,x1) -> (y1,y2) // permute params (PP)
	(x1,x2) -> (y1,y2) <~~> (x1,x2) -> (y2,y1) // permute results (PR)
2. 
	(x1,x2) -> (y1,y2) ~~> (x1,x2) -> y2 // weaken results (WR)
3.  // Motivation: query `func(string)int` => get too many `func() int` before the expected `strconv.Atoi :: func(string) (int, error)`  
    (x1,x2) -> (y1,y2) <~~ x2 -> (y1,y2) // weaken params (WP)
    
```

```
- subtyping

- currying? convenient or trouble?

- a -> b ~~ a -> c // possible first step
- a -> b ~~ c -> b // possible last step
```

```
- Fuzzy context: “Type name, type only name” ~~> Guess what user means (module, package, version)
    - type (v.)
    - according to popularity?

- Beyond expectation:
    - `func([]string) []int` ~~> `func[a,b any]([]a) []b`
    - `func([]string) []int` ~~> `func[a,b any]([]a, func(a) b) []b` + `func(string) int`

...
```
 
## Benchmark

```shell
cd ranking
go test -v . -test.bench ^BenchmarkDistance -test.run ^BenchmarkDistance -benchtime=200x -benchmem
```