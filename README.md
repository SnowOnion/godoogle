# Godoogle

[Hoogle](https://hoogle.haskell.org/) for [Go](https://go.dev/).

## Usage

### WIP: The Online One 

###  Deploy Your Own Godoogle

```shell
cd server/
go run *.go
```
Then visit [localhost:8888](http://localhost:8888).

WIP: configurable

## Challenges & Ideas

- Fuzzy context: “Type name, type only name” ~~> Guess what user means (module, package, version)
    - type (v.)
    - according to popularity?

- Beyond expectation:
    - `func([]string) []int` ~~> `func[a,b any]([]a) []b`
    - `func([]string) []int` ~~> `func[a,b any]([]a, func(a) b) []b` + `func(string) int`
    
...

## Approaches & TODOs

### Engineering

-[x] https://pkg.go.dev/std
-[x] https://github.com/samber/lo
-[ ] Many more...
  - [ ] Cover https://pkg.go.dev ? (Not the functionality, just the range)
-[ ] User specified modules / repos

### Algorithm

Learning from Hoogle...