# Godoogle

Godoogle is a [Go](https://go.dev/) API search engine, which (by now) allows you to search by **function type signature**, including [generic](https://go.dev/doc/tutorial/generics) functions.

Aim: Godoogle is to Go what [Hoogle](https://hoogle.haskell.org/) is to [Haskell](https://www.haskell.org/).

## Usages

### 🌏The Online One

https://godoogle.sonion.xyz/

### 🏡Deploy Your Own Godoogle

```shell
# per version
cd script/
go run floydWrite.go
 

# per run
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

- [ ] Fuzzy search by distance
- [ ] Learning from Hoogle...

## Challenges & Ideas

See [draft.md](draft.md).

## Licence

Not specified. You may suggest one!