# Godoogle

Godoogle is a [Go](https://go.dev/) API search engine, which allows you to search by approximate **function type signature**, including [generics](https://go.dev/doc/tutorial/generics).

Aim: Godoogle is to Go what [Hoogle](https://hoogle.haskell.org/) is to [Haskell](https://www.haskell.org/).

## Usages

### 🌏The Online One

https://godoogle.sonion.xyz/

### 🏡Deploy Your Own Godoogle

```shell
# First run
go run script/floydWrite.go
mv sigGraph.json server/res/

# Each run
cd server/
go run .
```

Then visit [localhost:8888](http://localhost:8888).

WIP: Be configurable. (By now: Mutate `func InitFuncDatabase()` in [collect/candidates.go](collect/candidates.go) ...)

## Approaches & TODOs

### Engineering

- [x] https://pkg.go.dev/std
- [x] https://github.com/samber/lo
- [ ] Many more...
    - [ ] Cover https://pkg.go.dev ? (Not the functionality, just the range)
- [ ] User specified modules / repos (i.e. configurable)

### Algorithm

- [x] Fuzzy search by distance
- [ ] Learning from Hoogle...

## Challenges & Ideas

See [draft.md](draft.md).

## Licence

Not specified. You may suggest one!