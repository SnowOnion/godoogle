# Godoogle

Godoogle is a [Go](https://go.dev/) API search engine, which allows you to search by approximate **function type
signature**, including [generics](https://go.dev/doc/tutorial/generics).

Aim: Godoogle is to Go what [Hoogle](https://hoogle.haskell.org/) is to [Haskell](https://www.haskell.org/).

## Usages

### 🌏The Online One

https://godoogle.sonion.xyz/

### 🏡Deploy Your Own Godoogle

```shell
# First run
go run script/writeSigGraph.go
mv sigGraph.json server/res/

# Each run
cd server/
go run .
```

Then visit [localhost:8888](http://localhost:8888).

WIP: Be configurable. (By now: Mutate `func InitFuncDatabase()` in [collect/candidates.go](collect/candidates.go) ...)

## Approaches & TODOs

- Smarter
    - [x] Fuzzy search by distance.
    - Adapt to various inputs.
        - [x] Type params, e.g. `[T any] func(bool, T, T) T`
        - [ ] Support non-[builtin](https://pkg.go.dev/builtin) types,
          e.g. `func() time.Time`, `func(l sync.Locker) *sync.Cond`
        - [ ] Copy straight from code / Godoogle result, e.g. `func InSlice[T comparable](item T, slice []T) bool`
        - [ ] No `func` nor name, e.g. `[T comparable](item T, slice []T) bool`
        - [ ] Omit type param, e.g. `func(...T) []T`, `func (T, []T) bool`
        - [ ] Package name with{,out} import path, e.g. `func(string) (*url.URL)` <-> `func(string) (*net/url.URL)`
    - Imagine more.
        - [x] `func(X,Y) (Z)` -> `func(X) (Z,W)`
        - [ ] `func(X,Y) (Z)` <- `func(X) (Z,W)`
        - [ ] `func(X,Y) (Z,W)` -> `func(Y,X) (W,Z)`
        - [ ] `[T any, F comparable]` <-> `[F comparable, T any]`
        - [ ] `[T comparable]func(...T) []T` -> `[T any]func(...T) []T`
        - [ ] `[E comparable](s []E, v E)` -> `[S ~[]E, E comparable](s S, v E)`
        - [ ] `func(X)` <-> `func(*X)`
    - First things first.
        - [ ] 
          Query `func(string) int` ⊢ `func(s string) (int, error)` > `func(s string) (p vendor/golang.org/x/text/unicode/bidi.Properties, sz int)`
    - [ ] Learn from Hoogle, [Roogle](https://roogle.hkmatsumoto.com/), *oogle.
- Wider
    - [x] [Standard library](https://pkg.go.dev/std) and a few 3rd party libs.
    - [ ] Support methods.
    - [ ] **Text-based search candidates, rather than “import-based”.**
    - [ ] Show https://pkg.go.dev/builtin#max (treated as not exported now).
    - [ ] Cover https://pkg.go.dev (Not the functionality, just the range).
    - [ ] Customize search candidates, for self-hosted user.
- Faster
    - Memoization of shortest paths in SigGraph, the graph with func signature as vertex and distance as edge weight.
        - [x] Text.
        - [ ] DB (RDB? Graph DB?).
- Other
    - [x] Google Analytics.
    - [ ] Google Analytics options, for self-hosted user.

## Challenges & Ideas

See [draft.md](draft.md).

## Licence

Not specified. You may suggest one!